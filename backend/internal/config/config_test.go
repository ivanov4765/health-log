package config

import (
	"strings"
	"testing"
	"time"
)

type mapSource map[string]string

func (source mapSource) LookupEnv(key string) (string, bool) {
	value, ok := source[key]
	return value, ok
}

func validSource() mapSource {
	return mapSource{
		"HEALTH_LOG_ENVIRONMENT":      "development",
		"AWS_REGION":                  "eu-west-3",
		"HEALTH_LOG_TABLE_NAME":       "HealthLogs-dev",
		"HEALTH_LOG_DASHBOARD_ORIGIN": "http://localhost:8080",
		"TELEGRAM_BOT_TOKEN":          "telegram-token-value",
		"TELEGRAM_WEBHOOK_SECRET":     "webhook-secret-value",
		"TELEGRAM_ALLOWED_USER_IDS":   "123, 456",
		"DASHBOARD_PASSWORD_HASH":     "$argon2id$placeholder",
		"SESSION_SIGNING_KEY":         "01234567890123456789012345678901",
	}
}

func TestLoadUsesDefaultsAndRedactsSecrets(t *testing.T) {
	config, err := Load(validSource())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if config.Environment != Development || config.SessionTTL != 8*time.Hour || config.SessionCookieName != "health_log_session" {
		t.Fatalf("Load() defaults did not match expected configuration")
	}
	if len(config.TelegramAllowedUserIDs) != 2 || config.TelegramAllowedUserIDs[1] != 456 {
		t.Fatalf("Load() user IDs = %#v", config.TelegramAllowedUserIDs)
	}
	if config.AIEnabled() {
		t.Fatal("AIEnabled() = true without Gemini configuration")
	}
	if got := config.TelegramBotToken.String(); got != "[REDACTED]" {
		t.Fatalf("Secret.String() = %q", got)
	}
}

func TestLoadRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		change func(mapSource)
		key    string
	}{
		{"missing table", func(values mapSource) { delete(values, "HEALTH_LOG_TABLE_NAME") }, "HEALTH_LOG_TABLE_NAME"},
		{"invalid environment", func(values mapSource) { values["HEALTH_LOG_ENVIRONMENT"] = "local" }, "HEALTH_LOG_ENVIRONMENT"},
		{"HTTP production origin", func(values mapSource) { values["HEALTH_LOG_ENVIRONMENT"] = "production" }, "HEALTH_LOG_DASHBOARD_ORIGIN"},
		{"origin with path", func(values mapSource) { values["HEALTH_LOG_DASHBOARD_ORIGIN"] = "https://example.com/app" }, "HEALTH_LOG_DASHBOARD_ORIGIN"},
		{"duplicate telegram IDs", func(values mapSource) { values["TELEGRAM_ALLOWED_USER_IDS"] = "123,123" }, "TELEGRAM_ALLOWED_USER_IDS"},
		{"invalid password hash", func(values mapSource) { values["DASHBOARD_PASSWORD_HASH"] = "not-a-hash" }, "DASHBOARD_PASSWORD_HASH"},
		{"short signing key", func(values mapSource) { values["SESSION_SIGNING_KEY"] = "short" }, "SESSION_SIGNING_KEY"},
		{"invalid TTL", func(values mapSource) { values["SESSION_TTL"] = "2h30" }, "SESSION_TTL"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values := validSource()
			test.change(values)
			_, err := Load(values)
			if err == nil || !strings.Contains(err.Error(), test.key) {
				t.Fatalf("Load() error = %v, want key %q", err, test.key)
			}
		})
	}
}

func TestLoadDoesNotEchoSecretInErrors(t *testing.T) {
	values := validSource()
	values["SESSION_SIGNING_KEY"] = "unsafe-secret"
	_, err := Load(values)
	if err == nil {
		t.Fatal("Load() error = nil")
	}
	if strings.Contains(err.Error(), "unsafe-secret") {
		t.Fatalf("Load() leaked secret in error: %v", err)
	}
}
