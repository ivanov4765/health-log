// Package events defines the health-event contract shared by capture and API
// handlers. Persistence and transport adapters are deliberately kept outside
// this package.
package events

import "time"

const (
	MaxNameLength       = 120
	MaxMedicationLength = 120
	MaxUnitLength       = 32
	MaxNoteLength       = 1000
	DefaultPageSize     = 50
	MaxPageSize         = 100
)

// Type describes a health event category.
type Type string

const (
	Symptom     Type = "symptom"
	Medication  Type = "medication"
	Measurement Type = "measurement"
)

// Event is the stored, server-owned representation of a health event.
type Event struct {
	ID         string
	OwnerID    string
	Type       Type
	OccurredAt time.Time
	Timezone   string
	Name       string
	Severity   *int
	Medication string
	Amount     *float64
	Unit       string
	Note       string
	Source     Source
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Source identifies how an event was captured.
type Source string

const (
	SourceTelegram  Source = "telegram"
	SourceDashboard Source = "dashboard"
)

// CreateInput contains caller-supplied event values. Owner identity and audit
// fields are intentionally absent because the server assigns them.
type CreateInput struct {
	Type       Type
	OccurredAt time.Time
	Timezone   string
	Name       string
	Severity   *int
	Medication string
	Amount     *float64
	Unit       string
	Note       string
	Source     Source
}

// UpdateInput carries mutable event fields. The API layer selects an existing
// owned event before applying the validated replacement values.
type UpdateInput = CreateInput

// PageRequest defines range-query pagination. Cursor is opaque to callers.
type PageRequest struct {
	From   time.Time
	To     time.Time
	Type   *Type
	Cursor string
	Limit  int
}

// Page is a range-query result.
type Page struct {
	Events     []Event
	NextCursor string
}

// Clock supports deterministic service tests.
type Clock interface {
	Now() time.Time
}

// SystemClock is the production clock.
type SystemClock struct{}

// Now returns the current UTC instant.
func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
