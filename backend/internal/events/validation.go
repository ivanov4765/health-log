package events

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

// ValidationError is safe to expose as a client error. It never includes the
// user's supplied value, which may itself be sensitive health information.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Message)
}

// Validate checks the public event contract before an event is persisted.
func (input CreateInput) Validate() error {
	if input.Type != Symptom && input.Type != Medication && input.Type != Measurement {
		return invalid("type", "must be symptom, medication, or measurement")
	}
	if input.OccurredAt.IsZero() {
		return invalid("occurredAt", "is required")
	}
	if strings.TrimSpace(input.Timezone) == "" {
		return invalid("timezone", "is required")
	}
	if !utf8.ValidString(input.Note) || utf8.RuneCountInString(input.Note) > MaxNoteLength {
		return invalid("note", "must be valid UTF-8 and at most 1000 characters")
	}
	if input.Source != SourceTelegram && input.Source != SourceDashboard {
		return invalid("source", "must be telegram or dashboard")
	}

	switch input.Type {
	case Symptom:
		if err := validateText("name", input.Name, MaxNameLength, true); err != nil {
			return err
		}
		if input.Severity == nil || *input.Severity < 0 || *input.Severity > 10 {
			return invalid("severity", "must be an integer from 0 through 10")
		}
		if input.Medication != "" || input.Amount != nil || input.Unit != "" {
			return invalid("symptom", "must not contain medication fields")
		}
	case Medication:
		if err := validateText("medication", input.Medication, MaxMedicationLength, true); err != nil {
			return err
		}
		if input.Name != "" || input.Severity != nil {
			return invalid("medication", "must not contain symptom fields")
		}
		if err := validateAmountAndUnit(input.Amount, input.Unit); err != nil {
			return err
		}
	case Measurement:
		if err := validateText("name", input.Name, MaxNameLength, true); err != nil {
			return err
		}
		if input.Medication != "" || input.Severity != nil {
			return invalid("measurement", "must not contain symptom or medication fields")
		}
		if input.Amount == nil || math.IsNaN(*input.Amount) || math.IsInf(*input.Amount, 0) {
			return invalid("amount", "must be a finite number")
		}
		if err := validateText("unit", input.Unit, MaxUnitLength, true); err != nil {
			return err
		}
	}

	return nil
}

// Validate checks the date-window, type filter, and page-size contract.
func (request PageRequest) Validate() error {
	if request.From.IsZero() || request.To.IsZero() {
		return invalid("dateRange", "from and to are required")
	}
	if !request.From.Before(request.To) {
		return invalid("dateRange", "from must be before to")
	}
	if request.To.Sub(request.From) > 366*24*60*60*1e9 {
		return invalid("dateRange", "must not exceed 366 days")
	}
	if request.Type != nil && *request.Type != Symptom && *request.Type != Medication && *request.Type != Measurement {
		return invalid("type", "must be symptom, medication, or measurement")
	}
	if request.Limit < 0 || request.Limit > MaxPageSize {
		return invalid("limit", "must be between 1 and 100")
	}
	return nil
}

// NormalizedLimit returns the default when callers omit a page size.
func (request PageRequest) NormalizedLimit() int {
	if request.Limit == 0 {
		return DefaultPageSize
	}
	return request.Limit
}

func validateAmountAndUnit(amount *float64, unit string) error {
	if amount == nil && strings.TrimSpace(unit) != "" {
		return invalid("amount", "is required when unit is provided")
	}
	if amount != nil && (math.IsNaN(*amount) || math.IsInf(*amount, 0) || *amount <= 0) {
		return invalid("amount", "must be a finite positive number")
	}
	if amount != nil {
		return validateText("unit", unit, MaxUnitLength, true)
	}
	return nil
}

func validateText(field, value string, maximum int, required bool) error {
	if !utf8.ValidString(value) {
		return invalid(field, "must be valid UTF-8")
	}
	trimmed := strings.TrimSpace(value)
	if required && trimmed == "" {
		return invalid(field, "is required")
	}
	if utf8.RuneCountInString(trimmed) > maximum {
		return invalid(field, fmt.Sprintf("must be at most %d characters", maximum))
	}
	return nil
}

func invalid(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}
