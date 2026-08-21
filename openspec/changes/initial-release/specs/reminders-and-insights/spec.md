# Reminders and insights specification

## Requirement: reminders are owner-controlled and reliable

The system SHALL let an authorized owner enable, change, and disable reminders.
It SHALL send enabled reminders through Telegram and record delivery failures
without sensitive text in logs.

### Scenario: disabled reminder does not send

- **WHEN** an owner disables a reminder before its next due occurrence
- **THEN** the next occurrence produces no Telegram notification

## Requirement: insights are bounded, descriptive, and optional

The system SHALL use Gemini only after an authenticated request. It SHALL send
only selected, bounded, owner-authorized events, exclude internal identifiers,
and frame output as non-diagnostic information with its context range shown.

### Scenario: insight over selected history

- **WHEN** a signed-in user asks for an insight over a valid date range
- **THEN** the result identifies range/event count, includes a disclaimer, and makes no diagnosis or prescription

### Scenario: Gemini is unavailable

- **WHEN** Gemini times out or returns an error
- **THEN** the system returns a recoverable unavailable message and normal logging/history continue to work
