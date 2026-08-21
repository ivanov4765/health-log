# Initial release: personal health log

**Status:** Revised for a single personal AWS account
**Owner:** One person
**Target:** v0.1.0  
**Source:** `project_proposa.md`, simplified for personal use

## Intent

Build a small private journal for one owner: log medication, symptoms, and
measurements from Telegram; view recent history in a static web page; and
optionally receive one daily reminder or request a Gemini summary. This is a
learning/personal project, not a multi-user product or production service.

## Definition of useful

The first usable version is complete when the owner can send a valid `/log`
message to Telegram, see the saved entry in the dashboard, and confirm that a
random visitor to the public dashboard cannot read the data without the access
token. All other features build on that vertical slice.

## Scope

### Included in v0.1.0

- One AWS account, one region, one deployed application environment, and a small
  S3-backed Terraform state. The state bucket and GitHub OIDC role are
  bootstrapped once.
- One Telegram bot, restricted to one configured Telegram user ID.
- One Lambda Function URL, one on-demand DynamoDB table, and one static GitHub
  Pages dashboard.
- Explicit `/log` and `/help` Telegram commands; no natural-language AI parser.
- Dashboard token prompt, last-30-days event list, simple charts, and CSV/JSON
  export. Dashboard editing and deletion are deferred.
- Optional Gemini insight over a selected range and one fixed daily Telegram
  reminder, after the logging/dashboard MVP works.
- Basic CloudWatch logs with short retention and a small AWS Budget or Free Tier
  alert set manually in the AWS console.
- GitHub Actions deploys the Lambda package and Terraform automatically on a
  push to `main`, using short-lived AWS OIDC credentials.

### Not included

- Multiple users, accounts, roles, sharing, caregiver access, or organization
  features.
- Separate development/staging/production environments, pull-request approval
  gates, complex release promotion, or uptime targets.
- Dashboard password/session flow, cookies, CSRF infrastructure, WAF, complex
  rate limiting, custom metrics/alarms, backups/PITR, and restore drills.
- Natural-language interpretation during capture, medication advice, diagnosis,
  emergency guidance, native apps, wearables, uploads, or full-text search.

## Minimal privacy boundary

Health data is still sensitive even for a personal project. The intentionally
small baseline is: keep all provider credentials and the dashboard access token
outside Git, store them in AWS SSM Parameter Store SecureString, allow only the
owner's Telegram ID, verify Telegram's webhook secret, and require a custom
dashboard token header for every data API. Public dashboard files contain no
health data or token.

This is a pragmatic personal-project boundary, not a compliance claim.

## Milestones

1. **Automation bootstrap:** owner creates the state bucket and GitHub OIDC role
   once; pushes to `main` then deploy Lambda and Terraform automatically.
2. **MVP logging:** Terraform deploys DynamoDB/Lambda; owner can log a single
   event through Telegram and retrieve it through a protected API.
3. **MVP dashboard:** GitHub Pages asks for the access token and shows recent
   events and a chart.
4. **Nice-to-haves:** add the fixed daily reminder and Gemini summary only if
   they remain useful after regular logging begins.
5. **Personal launch:** enable the bot webhook, set a small cost alert, run the
   manual smoke checklist, and use the project.

## Small risk register

| Risk | Practical response |
| --- | --- |
| Token or bot credential is committed | Keep `.env` ignored, use SSM for deployed secrets, rotate the affected secret. |
| Automatic deploy gets AWS access | Limit the OIDC role to this GitHub repository and `main`; use no stored AWS access keys. |
| A stranger sends bot updates | Check Telegram secret header and hard-code/configure the single allowed Telegram ID. |
| Public Pages site exposes data | Do not embed data/token in assets; API rejects requests without `X-Health-Log-Token`. |
| Free-tier usage grows unexpectedly | Keep one region/resource set, use on-demand DynamoDB, set a small manual budget alert, and remove unused resources. |
| Gemini is unhelpful or unavailable | Leave it disabled; core logging does not depend on it. |

## Release check

Before personal use, complete the MVP tasks, confirm a `main` push completes the
deployment workflow, run the manual smoke checklist with non-sensitive test
entries, and verify the AWS budget/Free Tier alert is enabled. There is no
staging or formal production release process.
