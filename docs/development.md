# Development setup

## Supported toolchain

The project pins its expected versions in `.tool-versions`:

- Go 1.26.3
- Node.js 23.4.0 (reserved for the dashboard toolchain)
- Terraform 1.13.3 (required beginning with the infrastructure phase)

An asdf-compatible version manager can install these directly. Go is the only
required tool for the current foundation. The Go module and service source live
in `backend/`. Do not install dependencies into the repository globally or
commit generated tool output.

## Configuration

Copy `backend/.env.example` to `backend/.env` and replace the placeholders. The
application does not load `.env` by itself; use your local process manager or
shell to export the variables. Production values must come from managed secret
storage, not a dotenv file.

Required variables are:

| Variable | Purpose |
| --- | --- |
| `HEALTH_LOG_ENVIRONMENT` | `development`, `staging`, or `production`. |
| `AWS_REGION` | Region that owns the health-log resources. |
| `HEALTH_LOG_TABLE_NAME` | DynamoDB table name. |
| `HEALTH_LOG_DASHBOARD_ORIGIN` | Exact dashboard origin. HTTP is allowed only in development. |
| `TELEGRAM_BOT_TOKEN` | Telegram API credential. |
| `TELEGRAM_WEBHOOK_SECRET` | Telegram webhook verification secret. |
| `TELEGRAM_ALLOWED_USER_IDS` | Comma-separated positive Telegram user IDs. |
| `DASHBOARD_PASSWORD_HASH` | Argon2id or bcrypt password verifier; never plaintext. |
| `SESSION_SIGNING_KEY` | At least 32 random characters for session signing. |

`GEMINI_API_KEY` is optional. Its absence disables AI insights. `SESSION_TTL`
defaults to `8h` and must be between 15 minutes and 24 hours.

## Commands

| Command | Purpose |
| --- | --- |
| `make -C backend fmt` | Format Go source. |
| `make -C backend lint` | Run Go static analysis. |
| `make -C backend test` | Run unit tests. |
| `make -C backend build` | Compile all Go packages. |
| `make -C backend audit` | Verify Go module integrity and show the dependency graph. |
| `make -C backend vulncheck` | Run a pinned Go vulnerability scan (requires network on first run). |
| `make -C backend check` | Run lint, tests, build, and dependency-integrity audit. |

For a dependency vulnerability scan, run the pinned command below from a
network-enabled development or CI environment:

```sh
make -C backend vulncheck
```

This command intentionally uses an explicit version so audits are reproducible.
It will be moved into required CI in OpenSpec task 7.1.
