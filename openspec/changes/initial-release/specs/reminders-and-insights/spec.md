# Optional reminders and insights specification

Neither capability is required for the logging/dashboard MVP.

## Requirement: one fixed daily reminder can be enabled

When enabled in Terraform, EventBridge SHALL send one configured daily Telegram
message to the owner. Changing or disabling it is done by changing Terraform
variables and applying them.

### Scenario: reminder is disabled

- **WHEN** the reminder variable is disabled and Terraform is applied
- **THEN** no future reminder is sent

## Requirement: Gemini insight is explicitly requested and informational

When Gemini configuration exists, the API MAY send at most 100 events from the
owner's selected date range to Gemini after a valid dashboard-token request.
The response SHALL include a short informational/not-medical-advice notice.

### Scenario: Gemini is not configured

- **WHEN** `GEMINI_API_KEY` is absent
- **THEN** the insight endpoint reports that insights are unavailable and event
  capture/history continue to work
