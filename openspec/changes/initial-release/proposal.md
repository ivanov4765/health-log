# Initial release: personal health event logging

**Status:** Draft  
**Owner:** Project maintainer  
**Target:** v0.1.0  
**Source:** `project_proposa.md`, version 1.0.0

## Intent

Deliver a low-maintenance personal health journal that lets an authorized person
capture medication, symptom, and measurement events through Telegram; review
trends in a static dashboard; receive reminders; and ask for descriptive,
non-diagnostic summaries of selected records.

## Outcomes and success measures

| Outcome | Release measure |
| --- | --- |
| Fast capture | A valid Telegram message receives confirmation within 5 seconds at p95, excluding provider outages. |
| Accurate history | A created event is visible in its date-range query; a Telegram retry cannot duplicate it. |
| Private access | Automated tests reject cross-user access and unauthenticated health-data reads/writes. |
| Useful history | Dashboard filters and visualizes medication adherence, symptom severity, and measurements. |
| Controlled AI | AI receives only selected authorized records, carries a disclaimer, and exposes no server secrets. |
| Low-cost operation | A fresh deployment has defined limits, a budget alert, and OIDC-only deployment. |

## In scope

- Go Lambdas, Terraform, static dashboard, tests, and GitHub Actions in one repository.
- Telegram capture of symptoms, medication intake, and numeric measurements.
- DynamoDB persistence, query API, user preferences, and idempotent webhooks.
- Configurable daily reminders through EventBridge Scheduler.
- Dashboard sign-in, history, charts, event editing, export, and deletion.
- Gemini-backed summaries/questions over a selected bounded date range.
- OIDC deployment, monitoring, restore procedure, and owner documentation.

## Explicitly out of scope for v0.1.0

- Clinical diagnosis, dose recommendations, emergency triage, or provider integrations.
- Shared journals, caregiver access, organizations, or payments.
- Native apps, offline sync, wearables, file uploads, and full-text search.
- Guaranteed zero cost, HIPAA certification, or a production SLA.
- AI-driven changes to medication schedules.

## Delivery assumptions requiring confirmation before production use

| Assumption | Initial-release decision |
| --- | --- |
| User model | One owner is supported end-to-end; data remains user-scoped for future isolation. |
| Locale | Instants are UTC; the dashboard stores an IANA timezone. Initial UI is English. |
| Retention | Events remain until owner deletion; backups follow the chosen policy. |
| Auth | Static assets are public, but every health-data API requires a password-derived short-lived session. CORS is not authorization. |
| Bot access | Only allowlisted Telegram IDs may use the bot. |
| AI provider | Gemini is optional; capture and history work when it is unavailable. |

## Milestones

1. **Foundation:** repository layout, data contract, threat model, and Terraform bootstrap are reviewed.
2. **Capture vertical slice:** an allowlisted Telegram user can create/retrieve an event in non-production.
3. **Dashboard vertical slice:** an authenticated user can browse, chart, create, update, export, and delete own events.
4. **Automation and intelligence:** reminders and bounded AI insights fail safely.
5. **Release readiness:** CI/CD, observability, restore drill, security tests, accessibility checks, and runbook are complete.

## Risks and mitigations

| Risk | Mitigation / release gate |
| --- | --- |
| Health-data exposure | Data minimization, encryption, managed secrets, redacted logs, allowlists, API auth, and authorization tests. |
| Duplicate webhook events | Conditional persistence of a Telegram update receipt. |
| Bypassed static-site password | No data in site bundle; session checked by every API. |
| Unsafe AI output | Bounded context, non-diagnostic rules, source-window metadata, and raw-record access. |
| Free-tier overrun | Limits, budgets, alerts, manual AI rate caps, and cost review. |
| Provider outage | Bounded retries and usable non-AI capture/history paths. |

## Release criteria

Create the v0.1.0 tag only after required `tasks.md` items are checked, specs
pass, staging smoke tests use synthetic data, and the security/cost checklist is
approved.
