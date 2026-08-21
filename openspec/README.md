# health-log OpenSpec

This directory is the delivery plan for a single-owner personal deployment of
`health-log`. It turns the proposal into small, verifiable steps without
pretending this is a multi-user production service.

## Change package

`changes/initial-release/` is the active change package.

- `proposal.md` defines scope, assumptions, success measures, and release boundaries.
- `design.md` records architecture and security decisions.
- `tasks.md` is the execution checklist; completion requires stated verification.
- `specs/` contains normative, testable capability requirements.

## Working agreement

1. Complete work in dependency order in `tasks.md` and keep evidence current.
2. Adjust the relevant `spec.md` before changing a user-visible behavior.
3. Never use real health data in tests, issues, screenshots, or AI prompts beyond the requested scope.
4. The $0 target is a budget guardrail, not a guarantee. Keep one resource set and enable a small AWS cost alert before use.
5. Keep the owner-only token, Telegram allowlist, webhook secret, and SSM-stored credentials even though the project is simple.
6. This is a personal-use tool, not a medical device. It does not diagnose, prescribe, or provide emergency guidance.

## Status lifecycle

`Draft` -> `In implementation` -> `Personal use`. A task is complete when its
focused automated check or listed manual check has been performed.
