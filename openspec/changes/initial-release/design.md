# Technical design: initial release

## Architecture

```text
Telegram -> Lambda Function URL (webhook) -> Go service -> DynamoDB
GitHub Pages -> Lambda Function URL (JSON API) -> Go service -> DynamoDB
                                                    \-> Gemini REST API (optional)
EventBridge Scheduler -> reminder Lambda -> Telegram sendMessage API
GitHub Actions -- OIDC --> AWS deployment role --> Terraform-managed resources
```

One Go `provided.al2023` Lambda binary may route webhook and API requests for
v0.1.0. The scheduler is separate and has only the permissions required to
dispatch reminders. This is an implementation boundary, not an API contract.

## Repository layout to implement

```text
cmd/api/                 Lambda bootstrap and HTTP router
internal/auth/           session and password verification
internal/events/         validation, service, repository interfaces
internal/telegram/       webhook validation, parser, Telegram client
internal/ai/             Gemini client, context selection, safety policy
internal/reminders/      schedule and dispatch logic
web/                     GitHub Pages HTML, CSS, JavaScript, assets
infra/terraform/         AWS resources, IAM, environments
tests/                   integration and end-to-end tests
.github/workflows/       validation and OIDC deployment pipelines
docs/                    setup, runbook, privacy notice
```

## Data model

Use one on-demand `HealthLogs` table. The service derives all keys; callers
cannot provide arbitrary ownership keys.

| Entity | PK | SK | Key fields |
| --- | --- | --- | --- |
| Event | `USER#{userId}` | `EVENT#{occurredAt}#{eventId}` | `type`, `occurredAt`, `timezone`, `value`, `unit`, `severity`, `medication`, `note`, `source`, audit times |
| User preference | `USER#{userId}` | `PROFILE` | Telegram ID, timezone, allow status, auth metadata |
| Reminder | `USER#{userId}` | `REMINDER#{reminderId}` | schedule, text, enabled state, next run |
| Webhook receipt | `USER#{userId}` | `UPDATE#{telegramUpdateId}` | received time, TTL expiration |

Event types are `symptom`, `medication`, and `measurement`. A symptom requires
name plus integer severity 0–10. A medication requires medication name and may
include amount/unit. A measurement requires name, numeric value, and unit.
Notes are optional and capped at 1,000 UTF-8 characters. All timestamps are RFC
3339 UTC instants; display uses the saved IANA timezone.

Start with a base-table range query per user; do not scan the table. Add an
index only when a measured access pattern cannot be served by this shape.

## API contract

All JSON API routes are under `/v1`, use `application/json`, return
`{ "error": { "code", "message" } }` on errors, and attach a request ID.
All except login and the Telegram webhook require a session.

| Method/path | Purpose |
| --- | --- |
| `POST /v1/auth/login` | Verify owner password and set a secure session. Rate limited. |
| `POST /v1/auth/logout` | Invalidate session. |
| `GET /v1/events?from=&to=&type=` | Paginated own-event range query; dates required and max 366 days. |
| `POST /v1/events` | Create validated own event. |
| `PATCH /v1/events/{id}` | Update own event only. |
| `DELETE /v1/events/{id}` | Delete own event only. |
| `GET /v1/export?from=&to=` | Download selected own events as JSON/CSV. |
| `GET/PUT /v1/preferences` | Read/update timezone and reminders. |
| `POST /v1/insights` | Gemini summary of selected dates/query. |
| `POST /telegram/webhook` | Telegram-only update receiver. |

Only configured GitHub Pages origins receive CORS headers. Authenticated browser
mutations require a signed per-session CSRF token. The webhook accepts no
browser origin.

## Authentication and secrets

- Store Telegram token, Gemini key, password verifier/pepper, session signing
  key, and webhook secret in Secrets Manager or SSM SecureString. Terraform
  receives references, never plaintext values in state.
- Validate `X-Telegram-Bot-Api-Secret-Token` with a constant-time comparison
  before parsing, then apply the Telegram user-ID allowlist.
- Store Argon2id or bcrypt hashes, never a dashboard password. Use short-lived,
  HttpOnly, Secure, SameSite sessions validated server-side on every request.
- Grant each Lambda only necessary DynamoDB actions/key scope and specific secret
  reads. No broad scans, wildcard secret access, or administrator roles.
- GitHub Actions assumes a role only through repository/branch/environment-bound
  OIDC conditions. No static AWS keys are used.

## AI safety and context rules

The client sends a question plus explicit dates. The API queries only the
authenticated user's records in that window, caps record count/bytes, redacts
internal IDs, and submits structured events to Gemini. The system prompt limits
output to descriptive patterns and uncertainty; it prohibits diagnosis,
prescribing, and emergency decisions. Results show the date window and number
of events used. Prompts and raw records are never logged. Gemini can be disabled
without affecting logging/history.

## Observability, resilience, and cost

- Structured logs contain request ID, route, result class, latency, and a hashed
  user ID only. Set short retention and redaction.
- Alarm on Lambda errors/throttles/duration, DynamoDB errors, webhook validation,
  reminder failure, AI failure, and budget threshold. Define an alert target.
- Retry Telegram/Gemini only with bounded exponential retries; never retry bad
  input. A conditional receipt write makes Telegram delivery idempotent.
- Set per-user/API rate limits plus event, date-range, and AI-context ceilings.
  Dashboard errors must be actionable and must not reveal internals.
- Before launch choose and document point-in-time recovery or encrypted periodic
  export/restoration, including cost impact, then test restoration with synthetic data.

## Owner configuration required before production

Set the GitHub Pages URL, Telegram IDs, alert destination, retention/log period,
backup option, monthly spend threshold, initial reminders, and secret
provisioning method. These are environment configuration, not source code.
