// Package config loads and validates runtime configuration without exposing
// secret values through error messages or String methods.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultCookieName = "health_log_session"
	defaultSessionTTL = 8 * time.Hour
	minSessionTTL     = 15 * time.Minute
	maxSessionTTL     = 24 * time.Hour
)

// Environment identifies the deployment tier.
type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

// Source makes configuration loading testable without mutating process state.
type Source interface {
	LookupEnv(key string) (string, bool)
}

type sourceFunc func(key string) (string, bool)

func (source sourceFunc) LookupEnv(key string) (string, bool) {
	return source(key)
}

// Secret keeps credential values out of ordinary formatted output.
type Secret struct {
	value string
}

// Reveal returns the secret for the narrow runtime call that requires it.
func (s Secret) Reveal() string {
	return s.value
}

// String deliberately redacts a secret in logs and errors.
func (s Secret) String() string {
	return "[REDACTED]"
}

// Config contains validated runtime configuration.
type Config struct {
	Environment            Environment
	AWSRegion              string
	TableName              string
	DashboardOrigin        string
	TelegramBotToken       Secret
	TelegramWebhookSecret  Secret
	TelegramAllowedUserIDs []int64
	DashboardPasswordHash  Secret
	SessionSigningKey      Secret
	GeminiAPIKey           Secret
	SessionCookieName      string
	SessionTTL             time.Duration
}

// AIEnabled reports whether the optional Gemini integration is configured.
func (c Config) AIEnabled() bool {
	return c.GeminiAPIKey.Reveal() != ""
}

// LoadFromOS loads configuration from the current process environment.
func LoadFromOS() (Config, error) {
	return Load(sourceFunc(os.LookupEnv))
}

// Load builds a Config from source. Error values name configuration keys but
// never include their values.
func Load(source Source) (Config, error) {
	environmentValue := optional(source, "HEALTH_LOG_ENVIRONMENT", string(Development))
	environment := Environment(environmentValue)
	if environment != Development && environment != Staging && environment != Production {
		return Config{}, invalid("HEALTH_LOG_ENVIRONMENT", "must be development, staging, or production")
	}

	awsRegion, err := required(source, "AWS_REGION")
	if err != nil {
		return Config{}, err
	}
	tableName, err := required(source, "HEALTH_LOG_TABLE_NAME")
	if err != nil {
		return Config{}, err
	}
	telegramBotToken, err := required(source, "TELEGRAM_BOT_TOKEN")
	if err != nil {
		return Config{}, err
	}
	telegramWebhookSecret, err := required(source, "TELEGRAM_WEBHOOK_SECRET")
	if err != nil {
		return Config{}, err
	}
	passwordHash, err := required(source, "DASHBOARD_PASSWORD_HASH")
	if err != nil {
		return Config{}, err
	}
	if !isPasswordHash(passwordHash) {
		return Config{}, invalid("DASHBOARD_PASSWORD_HASH", "must be an Argon2id or bcrypt verifier")
	}

	dashboardOriginRaw, err := required(source, "HEALTH_LOG_DASHBOARD_ORIGIN")
	if err != nil {
		return Config{}, err
	}
	dashboardOrigin, err := parseOrigin(dashboardOriginRaw, environment)
	if err != nil {
		return Config{}, err
	}

	allowedUserIDsRaw, err := required(source, "TELEGRAM_ALLOWED_USER_IDS")
	if err != nil {
		return Config{}, err
	}
	allowedUserIDs, err := parseUserIDs(allowedUserIDsRaw)
	if err != nil {
		return Config{}, err
	}

	sessionTTL, err := parseSessionTTL(optional(source, "SESSION_TTL", defaultSessionTTL.String()))
	if err != nil {
		return Config{}, err
	}

	signingKey, err := required(source, "SESSION_SIGNING_KEY")
	if err != nil {
		return Config{}, err
	}
	if len(signingKey) < 32 {
		return Config{}, invalid("SESSION_SIGNING_KEY", "must contain at least 32 characters")
	}

	return Config{
		Environment:            environment,
		AWSRegion:              awsRegion,
		TableName:              tableName,
		DashboardOrigin:        dashboardOrigin,
		TelegramBotToken:       Secret{value: telegramBotToken},
		TelegramWebhookSecret:  Secret{value: telegramWebhookSecret},
		TelegramAllowedUserIDs: allowedUserIDs,
		DashboardPasswordHash:  Secret{value: passwordHash},
		SessionSigningKey:      Secret{value: signingKey},
		GeminiAPIKey:           Secret{value: optional(source, "GEMINI_API_KEY", "")},
		SessionCookieName:      optional(source, "SESSION_COOKIE_NAME", defaultCookieName),
		SessionTTL:             sessionTTL,
	}, nil
}

func required(source Source, key string) (string, error) {
	value, ok := source.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return "", invalid(key, "is required")
	}
	return strings.TrimSpace(value), nil
}

func optional(source Source, key, fallback string) string {
	value, ok := source.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func parseOrigin(raw string, environment Environment) (string, error) {
	if raw == "" {
		return "", invalid("HEALTH_LOG_DASHBOARD_ORIGIN", "is required")
	}
	origin, err := url.ParseRequestURI(raw)
	if err != nil || origin.Host == "" || (origin.Scheme != "https" && origin.Scheme != "http") || origin.Path != "" || origin.RawQuery != "" || origin.Fragment != "" {
		return "", invalid("HEALTH_LOG_DASHBOARD_ORIGIN", "must be an origin without path, query, or fragment")
	}
	if environment != Development && origin.Scheme != "https" {
		return "", invalid("HEALTH_LOG_DASHBOARD_ORIGIN", "must use HTTPS outside development")
	}
	return origin.String(), nil
}

func parseUserIDs(raw string) ([]int64, error) {
	if raw == "" {
		return nil, invalid("TELEGRAM_ALLOWED_USER_IDS", "is required")
	}

	seen := make(map[int64]struct{})
	ids := make([]int64, 0)
	for _, item := range strings.Split(raw, ",") {
		id, err := strconv.ParseInt(strings.TrimSpace(item), 10, 64)
		if err != nil || id <= 0 {
			return nil, invalid("TELEGRAM_ALLOWED_USER_IDS", "must contain unique positive integer IDs")
		}
		if _, exists := seen[id]; exists {
			return nil, invalid("TELEGRAM_ALLOWED_USER_IDS", "must not contain duplicate IDs")
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

func parseSessionTTL(raw string) (time.Duration, error) {
	ttl, err := time.ParseDuration(raw)
	if err != nil || ttl < minSessionTTL || ttl > maxSessionTTL {
		return 0, invalid("SESSION_TTL", "must be between 15m and 24h")
	}
	return ttl, nil
}

func isPasswordHash(value string) bool {
	return strings.HasPrefix(value, "$argon2id$") ||
		strings.HasPrefix(value, "$2a$") ||
		strings.HasPrefix(value, "$2b$") ||
		strings.HasPrefix(value, "$2y$")
}

func invalid(key, reason string) error {
	return fmt.Errorf("invalid configuration for %s: %s", key, reason)
}
