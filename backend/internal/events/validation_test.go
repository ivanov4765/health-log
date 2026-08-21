package events

import (
	"strings"
	"testing"
	"time"
)

func intPointer(value int) *int           { return &value }
func floatPointer(value float64) *float64 { return &value }

func validSymptom() CreateInput {
	return CreateInput{
		Type: Symptom, OccurredAt: time.Date(2026, 8, 21, 9, 0, 0, 0, time.UTC),
		Timezone: "Europe/Paris", Name: "Headache", Severity: intPointer(4), Source: SourceTelegram,
	}
}

func TestCreateInputValidateAcceptsEachEventType(t *testing.T) {
	tests := []CreateInput{
		validSymptom(),
		{Type: Medication, OccurredAt: time.Now(), Timezone: "UTC", Medication: "Ibuprofen", Amount: floatPointer(200), Unit: "mg", Source: SourceDashboard},
		{Type: Measurement, OccurredAt: time.Now(), Timezone: "UTC", Name: "Blood pressure", Amount: floatPointer(120), Unit: "mmHg", Source: SourceDashboard},
	}
	for _, input := range tests {
		if err := input.Validate(); err != nil {
			t.Fatalf("Validate() error = %v for %#v", err, input.Type)
		}
	}
}

func TestCreateInputValidateRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		input func() CreateInput
		field string
	}{
		{"unknown type", func() CreateInput { value := validSymptom(); value.Type = "other"; return value }, "type"},
		{"missing timestamp", func() CreateInput { value := validSymptom(); value.OccurredAt = time.Time{}; return value }, "occurredAt"},
		{"missing timezone", func() CreateInput { value := validSymptom(); value.Timezone = ""; return value }, "timezone"},
		{"invalid severity", func() CreateInput { value := validSymptom(); value.Severity = intPointer(11); return value }, "severity"},
		{"medication has symptom fields", func() CreateInput {
			value := validSymptom()
			value.Type = Medication
			value.Medication = "Ibuprofen"
			return value
		}, "medication"},
		{"measurement requires amount", func() CreateInput {
			return CreateInput{Type: Measurement, OccurredAt: time.Now(), Timezone: "UTC", Name: "Heart rate", Unit: "bpm", Source: SourceDashboard}
		}, "amount"},
		{"note too long", func() CreateInput {
			value := validSymptom()
			value.Note = strings.Repeat("a", MaxNoteLength+1)
			return value
		}, "note"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.input().Validate()
			validationError, ok := err.(*ValidationError)
			if !ok || validationError.Field != test.field {
				t.Fatalf("Validate() error = %#v, want field %q", err, test.field)
			}
		})
	}
}

func TestPageRequestValidateAndNormalize(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	request := PageRequest{From: from, To: from.AddDate(0, 1, 0)}
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if got := request.NormalizedLimit(); got != DefaultPageSize {
		t.Fatalf("NormalizedLimit() = %d, want %d", got, DefaultPageSize)
	}

	request.To = from.AddDate(1, 1, 0)
	if err := request.Validate(); err == nil {
		t.Fatal("Validate() error = nil for range longer than 366 days")
	}
}
