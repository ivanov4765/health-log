# health-log

A serverless, personal health and medication logging platform built around
Telegram, AWS Lambda, DynamoDB, and a static dashboard.

The Go service is isolated in `backend/`; the static dashboard and Terraform
remain top-level workspaces. Product scope, delivery decisions, and the
implementation backlog live in the [OpenSpec package](openspec/README.md).

## Local development

Install the versions in `.tool-versions`, copy `backend/.env.example` to
`backend/.env`, and replace all placeholders locally. Do not commit `.env` or
any credentials.

```sh
make -C backend check
```

See [development documentation](docs/development.md) for individual commands,
configuration variables, and the dependency-vulnerability check.
