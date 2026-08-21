# Technical design: simple personal deployment

## Architecture

```text
Telegram -> Lambda Function URL -> Go Lambda -> DynamoDB
GitHub Pages -> Lambda Function URL -> Go Lambda -> DynamoDB
                                      \-> Gemini (optional)
EventBridge Scheduler -> Lambda -> Telegram (optional fixed daily reminder)
GitHub Actions -- OIDC --> Terraform + Lambda package -> one AWS account / one region
Owner laptop + AWS profile -> bootstrap state bucket and OIDC role (once)
```

Use one Go Lambda for the webhook and JSON API. A second small Lambda is needed
only if the fixed daily reminder is enabled. There is no API Gateway, Cognito,
queue, separate environment, or deployment-approval process.

## Repository layout

```text
backend/                 Go module: Lambda handler, DynamoDB, Telegram, Gemini
web/                     static HTML/CSS/JavaScript dashboard
infra/bootstrap/         one-time state bucket and GitHub OIDC deploy role
infra/terraform/         Lambda, table, URL, IAM; state stored in bootstrap bucket
.github/workflows/       automatic Lambda/Terraform and Pages deployment
tests/                   optional integration fixtures; unit tests stay with code
docs/                    short setup and personal runbook
openspec/                this plan
```

## Data model

Use one on-demand `HealthLogs` table with one fixed partition key. The owner is
identified by the configured Telegram ID or dashboard token, not a request field.

| Item | PK | SK | Fields |
| --- | --- | --- | --- |
| Health event | `OWNER` | `EVENT#{occurredAt}#{id}` | type, time, timezone, name/medication, severity/value/unit, note, source |
| Telegram receipt | `OWNER` | `UPDATE#{updateId}` | TTL expiry; prevents duplicate webhook delivery |

Keep the current three event types: `symptom`, `medication`, and `measurement`.
All event times are stored in UTC. Query the base table by date range; do not
add indexes or scans for v0.1.0.

## Minimal HTTP surface

| Route | Purpose |
| --- | --- |
| `POST /telegram/webhook` | Accept owner-only Telegram updates. |
| `GET /v1/events?from=&to=` | Return recent events; default to last 30 days. |
| `GET /v1/export?from=&to=` | Return selected records as JSON or CSV. |
| `POST /v1/insights` | Optional Gemini summary of selected records. |

There is no login route, session store, cookie, CSRF token, event editor, or
preferences API in v0.1.0. Event creation happens in Telegram. Keep responses
as straightforward JSON and return a short error message for invalid input.

## Dashboard access

The static dashboard asks the owner for a long random access token. JavaScript
sends it in `X-Health-Log-Token` on each API request and keeps it only in
`sessionStorage`. Lambda compares it to `DASHBOARD_ACCESS_TOKEN` in constant
time. The token is not committed or embedded in the page. The API sends CORS
headers only for the configured GitHub Pages origin. This custom header avoids
cookie-based CSRF complexity.

This is deliberately modest authentication suitable only for the sole owner.
Rotate the token in SSM and re-enter it in the dashboard if it is lost or
exposed.

## Secrets, state, and deployment

Store `TELEGRAM_BOT_TOKEN`, `TELEGRAM_WEBHOOK_SECRET`,
`DASHBOARD_ACCESS_TOKEN`, and optional `GEMINI_API_KEY` as SSM Parameter Store
SecureString values. Terraform references parameter names; it does not contain
their values.

The owner runs `infra/bootstrap` locally once to create a private versioned S3
state bucket, GitHub's OIDC provider, and one deploy role. The role trust policy
allows only this repository's `main` branch. The normal Terraform root uses the
S3 backend with encryption and lockfile locking. Store the bootstrap state
locally and ignore it; it is only needed to change the bootstrap resources.

GitHub Actions runs on pushes to `main` and manual dispatch. It checks the Go
module, builds a Linux Lambda `bootstrap` artifact, initializes Terraform using
the S3 state bucket, then plans and applies. It obtains short-lived credentials
with `id-token: write`; no AWS keys are stored in GitHub. Put the region, role
ARN, and state-bucket name in GitHub Actions variables. A separate Pages job can
publish `web/` when dashboard files change.

The Lambda role can read only those named parameters and access the one DynamoDB
table. Lambda default encryption at rest, HTTPS, and short CloudWatch log
retention are sufficient for this personal deployment.

## Capture, reminder, and AI behavior

- Start with deterministic `/log` syntax and `/help`. Reject malformed values
  with an example; do not use Gemini to guess what a message means.
- Verify Telegram's `X-Telegram-Bot-Api-Secret-Token` and the configured owner
  ID. Save update IDs conditionally so a retry cannot create two entries.
- If enabled, the reminder is one Terraform-configured daily time/message. It
  has no dashboard settings or schedule editor.
- If Gemini is configured, fetch only the requested date range, cap it to the
  most recent 100 events, and state that output is informational—not medical
  advice. Do not log the prompt or raw records. If it fails, show a simple
  unavailable message and continue logging normally.

## Operations

Use the AWS console's Budget or Free Tier alert with a small threshold. Inspect
CloudWatch and the GitHub Actions run only when debugging. The personal runbook
needs just: bootstrap, deploy/update, change secrets, set webhook, export
records, and delete Terraform resources when finished.
