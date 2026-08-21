# Dashboard access specification

## Requirement: health-data APIs require server-enforced authentication

The system SHALL serve public static dashboard assets without health data or
secrets. Every health-data API SHALL require a valid short-lived,
server-validated session.

### Scenario: unauthenticated data request

- **WHEN** a browser or script requests events, export, preferences, or insights without a valid session
- **THEN** the API rejects it and returns no health data

### Scenario: session expiry

- **WHEN** a dashboard session expires
- **THEN** protected requests are rejected and the dashboard guides a sign-in without discarding unsaved local form input

## Requirement: browser mutations resist cross-site requests

The system SHALL restrict CORS to configured dashboard origins and SHALL verify
a per-session CSRF token for authenticated browser mutations.

### Scenario: cross-site deletion attempt

- **WHEN** a third-party origin submits a cookie-bearing deletion without a valid CSRF token
- **THEN** the API rejects it and the event remains unchanged

## Requirement: dashboard presents user-controlled history

The dashboard SHALL allow date selection, event viewing, own-event CRUD, export
of a selected range, and charts whose labels/units match stored values.

### Scenario: user deletes an event

- **WHEN** the user confirms deletion of one event
- **THEN** only that event is removed from later reads and unrelated records remain intact
