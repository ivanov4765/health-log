# Execution tasks: personal MVP

This is a deliberately small one-account checklist. `[x]` means completed;
`[ ]` means pending. Verification is either a focused unit test or a recorded
manual check—not a production-release gate.

## 0. One-time owner setup

- [ ] **0.1 Record non-secret settings** — Write region, GitHub Pages URL,
  Telegram owner ID, timezone, and the desired fixed reminder time in
  `docs/personal-setup.md`. **Verify:** document contains no credential.
- [ ] **0.2 Create required secrets** — Put Telegram bot token, webhook secret,
  and random dashboard access token into named SSM SecureString parameters; add
  Gemini key only if using insights. **Verify:** `aws ssm get-parameter
  --with-decryption` works locally; no value is committed.
- [ ] **0.3 Enable a small cost alert** — Create an AWS Budget or Free Tier
  alert in the console. **Verify:** alert and email destination are visible in
  the billing console.
- [ ] **0.4 Bootstrap CI access once** — From the owner's local AWS profile,
  create a private versioned S3 Terraform-state bucket, GitHub OIDC provider,
  and deploy role restricted to this repository's `main` branch. Put region,
  role ARN, and state bucket in GitHub Actions variables. **Verify:** a small
  GitHub Actions OIDC test can assume the role without stored AWS keys.

## 1. Foundation

- [x] **1.1 Repository and Go module** — Go code lives in `backend/`; dashboard,
  Terraform, docs, and tests have their own top-level directories. **Verify:**
  `make -C backend build` succeeds.
- [x] **1.2 Local checks** — Tool versions, format/lint/test/build commands, and
  sample configuration are documented. **Verify:** `make -C backend check` succeeds.
- [x] **1.3 Event contract** — Event types, field validation, date-range query,
  pagination, and a clock abstraction exist with unit tests. **Verify:**
  `make -C backend test` succeeds.
- [ ] **1.4 Simplify runtime configuration** — Replace session/password settings
  with `DASHBOARD_ACCESS_TOKEN`; keep Gemini optional and require the owner
  Telegram ID. **Verify:** config tests reject missing token/owner ID and do
  not print secret values.

## 2. Automated AWS deployment

- [ ] **2.1 Write single-environment Terraform** — Create one table, one Lambda
  Function URL, log retention, IAM role, SSM parameter-name references, and a
  configured dashboard CORS origin. Configure the S3 state backend with
  encryption/lockfile locking. **Verify:** `terraform fmt`, `validate`, and
  `plan` succeed from `infra/terraform`.
- [ ] **2.2 Add the deploy workflow** — On `main` push and manual dispatch, run
  backend checks, build a Linux Lambda `bootstrap` ZIP, assume the GitHub OIDC
  role, then run Terraform init/plan/apply. **Verify:** workflow uses
  `id-token: write` and contains no AWS access-key secret.
- [ ] **2.3 Verify automatic deployment** — Merge a harmless Terraform or Lambda
  change to `main` and confirm the workflow updates the account and reports the
  Function URL. **Verify:** a simple health request reaches Lambda after CI runs.

## 3. Logging API and Telegram MVP

- [ ] **3.1 Store and query events** — Implement DynamoDB conditional event
  writes, update-receipt writes, and date-range reads for `OWNER`. **Verify:**
  unit/integration test creates, reads, and rejects a duplicate update ID.
- [ ] **3.2 Protect dashboard routes** — Check `X-Health-Log-Token` in constant
  time for event/export/insight routes; set the one allowed CORS origin.
  **Verify:** curl without/wrong token gets no data; correct token gets test data.
- [ ] **3.3 Implement Telegram commands** — Add `/help` and explicit `/log`
  parsing; verify webhook secret and owner Telegram ID. **Verify:** fixture tests
  reject bad secret, wrong owner, and malformed log.
- [ ] **3.4 Connect the real bot** — Set the Telegram webhook and send one
  non-sensitive test log. **Verify:** exactly one DynamoDB item and confirmation
  reply appear.

## 4. Dashboard MVP

- [ ] **4.1 Build one static page** — Add token prompt, last-30-days event list,
  simple Chart.js chart, loading/error state, and export link. **Verify:** open
  locally with synthetic data and confirm no token/data is in built assets.
- [ ] **4.2 Publish GitHub Pages** — Enable Pages once, then deploy `web/` from
  a GitHub Actions Pages job on `main`; set the deployed origin in Terraform.
  **Verify:** owner can enter token, see the test event, and export it;
  unauthenticated browser request fails.

## 5. Optional additions after MVP

- [ ] **5.1 Fixed daily reminder** — Add one EventBridge schedule configured by
  Terraform variables and a small Telegram sender. **Verify:** test schedule
  sends one reminder; changing or disabling it requires Terraform apply.
- [ ] **5.2 Simple Gemini insight** — Add a single dashboard question box and
  Lambda endpoint that submits at most 100 selected events. **Verify:** answer
  shows range/disclaimer; missing Gemini key leaves dashboard working.

## 6. Personal launch

- [ ] **6.1 Write a short personal runbook** — Document bootstrap, automatic
  deploy, secret/token rotation, webhook setup, export, and `terraform destroy`.
  **Verify:** it can be followed from a clean terminal session.
- [ ] **6.2 Run the manual smoke check** — Log a medication and symptom, view
  dashboard history, try incorrect token, export data, and check the cost
  alert. **Verify:** record results in `docs/personal-setup.md`.
