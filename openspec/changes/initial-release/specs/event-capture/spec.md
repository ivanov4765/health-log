# Event capture specification

## Requirement: only the owner can create a valid event through Telegram

The webhook SHALL verify Telegram's secret header and the one configured owner
Telegram ID before it processes an update. A valid `/log` command SHALL create
one validated symptom, medication, or measurement event with a UTC timestamp.

### Scenario: owner logs medication

- **WHEN** the configured Telegram owner sends a valid medication `/log`
- **THEN** the system saves one event and replies with a short confirmation

### Scenario: another user messages the bot

- **WHEN** a Telegram update comes from any user other than the configured owner
- **THEN** the system creates no event and returns no health data

### Scenario: Telegram retries an update

- **WHEN** the same Telegram update ID is delivered again
- **THEN** the system does not create a second event

## Requirement: recent history is simple to query

The system SHALL query the single owner's events by date range using the base
table and SHALL not scan the table.

### Scenario: recent dashboard history

- **WHEN** the dashboard requests a valid date range
- **THEN** it receives the owner's events in time order for that range only
