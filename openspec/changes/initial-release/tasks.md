# Execution tasks: initial release

**Status:** `[ ]` not started, `[-]` in progress, `[x]` complete.  
**Gate:** a required phase exit; do not start dependent work until its evidence exists.

## 0. Product and safety baseline

- [ ] **0.1 Confirm production settings** — Record owner, Pages URL, Telegram allowlist, timezone, budget threshold, alert route, backup choice, retention, and reminder copy in `docs/operations/configuration.md`. **Verify:** review confirms no secrets are recorded.
- [ ] **0.2 Publish privacy and safety copy** — Write retention, export/deletion, AI disclosure, non-medical, and emergency disclaimer copy. **Verify:** dashboard and README link it.
- [ ] **0.3 Threat model (Gate)** — Map assets, actors, trust boundaries, bot takeover, session theft, IDOR, spoofed webhook, prompt injection, and CI abuse to mitigations/tests. **Verify:** each high risk has an owner and test task.

## 1. Foundation and local development

- [ ] **1.1 Establish repository structure** — Add Go module, `cmd`, `internal`, `web`, `infra/terraform`, `tests`, `docs`, `.gitignore`, and placeholder-only `.env.example`. **Verify:** clean clone builds; no IDE/secrets/build output is tracked.
- [ ] **1.2 Pin tooling and checks** — Specify Go/Terraform/browser tooling plus format, lint, unit-test, and dependency-audit commands. **Verify:** one local command runs the full validation suite.
- [ ] **1.3 Implement domain contract** — Add event DTOs, validation, pagination, clock abstraction, and error codes. **Verify:** table-driven tests cover boundaries and invalid input.
- [ ] **1.4 Implement configuration loader** — Validate runtime settings at cold start and treat Gemini as optional. **Verify:** missing/malformed settings fail safely without secrets in output.

## 2. Infrastructure and deployment boundary

- [ ] **2.1 Terraform state and environments** — Create encrypted/locked remote state and isolated `dev`, `staging`, `prod` variables/workspaces. **Verify:** `fmt`, `validate`, and a non-prod plan succeed; state has no plaintext secret.
- [ ] **2.2 Provision compute and data** — Define DynamoDB, Functions/URLs, IAM, secret references, log retention, and origin-specific CORS. **Verify:** plan demonstrates least privilege and no public datastore endpoint.
- [ ] **2.3 Provision scheduler and alarms** — Define EventBridge schedules, retry/DLQ policy, CloudWatch metrics/alarms, and spend alert. **Verify:** controlled failure reaches test alert path.
- [ ] **2.4 Configure GitHub OIDC (Gate)** — Scope deploy-role trust to repository, branch/tag, and GitHub Environment; require production approval. **Verify:** workflow deploys using OIDC with no static AWS credentials.

## 3. Core persistence and protected API

- [ ] **3.1 DynamoDB repository** — Implement conditional CRUD, own-user range query, cursor pagination, profile/reminder/receipt access. **Verify:** DynamoDB Local tests prove no scan and no cross-user data access.
- [ ] **3.2 Session authentication** — Implement password verification, login throttle, secure short-lived sessions, logout, and CSRF. **Verify:** no/expired/tampered sessions and CSRF-free mutations fail.
- [ ] **3.3 Event and preference API** — Implement `/v1/events` and `/v1/preferences` handlers with ownership and range/schema validation. **Verify:** HTTP tests cover CRUD, pagination, malformed requests, and IDOR.
- [ ] **3.4 Export and deletion** — Stream JSON/CSV export and owned-event deletion. **Verify:** export contains only selected user/date range; deletion is immediately absent from reads.
- [ ] **3.5 API vertical slice (Gate)** — Deploy staging and exercise synthetic data. **Verify:** deployed smoke script passes and CORS/security headers are inspected.

## 4. Telegram capture

- [ ] **4.1 Webhook validator and client** — Constant-time secret check, update parsing, allowlist, and time-bounded Telegram reply client. **Verify:** forged-secret/non-allowlisted fixtures create no event.
- [ ] **4.2 Deterministic commands** — Implement `/help`, `/log`, `/stats`, syntax guidance, and errors. **Verify:** fixtures cover valid, malformed, unknown command, and timezone cases.
- [ ] **4.3 Assisted free text** — Begin with key-value parsing; any AI parsing must show explicit confirmation before persistence. **Verify:** ambiguous input cannot silently create/alter an event.
- [ ] **4.4 Idempotency and webhook registration** — Conditional receipt write before processing; configure HTTPS webhook secret. **Verify:** replayed update ID produces exactly one event.
- [ ] **4.5 Capture vertical slice (Gate)** — Run staging test with an allowlisted account. **Verify:** confirmation latency and persisted record meet success measures.

## 5. Dashboard

- [ ] **5.1 Accessible shell** — Build sign-in, expiry, responsive navigation, loading/error state, and privacy/safety links. **Verify:** keyboard-only journey and automated accessibility scan have no critical issue.
- [ ] **5.2 History and editor** — Add date/type filters, pagination, create/edit/delete confirmation, validation, and export. **Verify:** browser E2E tests complete every flow.
- [ ] **5.3 Charts** — Implement Chart.js symptom, adherence, and measurement views with units/no-data states. **Verify:** fixtures render correct labels/values; assets include no data.
- [ ] **5.4 Staging Pages deployment (Gate)** — Configure build, base path, CSP, HTTPS references, and API origin. **Verify:** deployed site authenticates and completes dashboard smoke flow.

## 6. Reminders and AI insights

- [ ] **6.1 Reminder lifecycle** — Validate preferences, create/update/disable schedules, dispatch idempotently, and record redacted delivery outcomes. **Verify:** one due reminder sends; disabled one does not; retry is observable.
- [ ] **6.2 Bounded Gemini client** — Implement timeout, context cap, record normalization/redaction, and error mapping. **Verify:** mocked provider receives only authorized selected-range data.
- [ ] **6.3 Insights API/UI** — Add question, date preview, failure state, context metadata, disclaimer, and raw-record link. **Verify:** Gemini failure leaves all non-AI functions usable.
- [ ] **6.4 AI safety review (Gate)** — Red-team advice, prompt injection, and data-exfiltration fixtures. **Verify:** approved safe responses become regression tests.

## 7. Release hardening and operations

- [ ] **7.1 CI/CD** — PR checks run format, lint, tests, web checks, Terraform plan, secret/dependency scans; approved main deploys staging; production needs approval. **Verify:** deliberate failing fixture blocks merge/deploy.
- [ ] **7.2 Performance and limits** — Measure cold/warm path, webhook latency, pagination, duplicate concurrency, rate limit, and payload limits. **Verify:** results meet or revise release measures with rationale.
- [ ] **7.3 Restore drill** — Restore/export synthetic data to isolation and query it. **Verify:** time/result/corrections recorded in runbook.
- [ ] **7.4 Operations runbook** — Cover setup, rotation, allowlist, outage, cost alarm, backup/restore, export/delete, teardown, and rollback. **Verify:** maintainer follows it in staging without source changes.
- [ ] **7.5 Release candidate (Gate)** — Tag RC, deploy staging, run specs/manual synthetic smoke tests, and review cost/security. **Verify:** evidence exists for each task; no unresolved critical/high issue.
- [ ] **7.6 Production release** — Approve production plan, deploy, register bot webhook, smoke test, and release notes. **Verify:** alarms are green and owner can capture/view/delete an event.
