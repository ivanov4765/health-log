# Event capture specification

## Requirement: authorized users can record valid health events

The system SHALL create an event only for an authorized user and only after
validating its type-specific fields. The server SHALL assign event identity,
ownership, and audit timestamps; it SHALL store timestamps as UTC instants.

### Scenario: medication intake is logged through Telegram

- **WHEN** an allowlisted Telegram user submits a valid medication log
- **THEN** the system creates exactly one medication event for that user and replies with normalized medication/time

### Scenario: invalid symptom severity is rejected

- **WHEN** a user submits a symptom with severity outside 0 through 10
- **THEN** no event is persisted and the user receives a corrective error

### Scenario: duplicate Telegram delivery is idempotent

- **WHEN** Telegram delivers the same update ID more than once
- **THEN** at most one event is created and duplicate handling creates no repeated side effect

## Requirement: records remain isolated by owner

The system SHALL derive its user partition key from authenticated server-side
identity and SHALL never trust a client-provided owner identifier.

### Scenario: a caller attempts cross-user access

- **WHEN** an authenticated user requests, updates, or deletes another user's event
- **THEN** the system returns a non-disclosing not-found/forbidden response and makes no change

## Requirement: history is bounded and queryable

The system SHALL return only the caller's events within a required bounded date
range and SHALL paginate results with an opaque cursor.

### Scenario: dashboard requests a filtered range

- **WHEN** a signed-in user requests events for a valid date range and type
- **THEN** only matching owned events are returned in stable time order, with a cursor if more exist
