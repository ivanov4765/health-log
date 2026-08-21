# health-log OpenSpec

This directory is the delivery contract for the initial public release of
`health-log`. It turns the proposal into independently verifiable requirements,
technical decisions, and an ordered implementation backlog.

## Change package

`changes/initial-release/` is the active change package.

- `proposal.md` defines scope, assumptions, success measures, and release boundaries.
- `design.md` records architecture and security decisions.
- `tasks.md` is the execution checklist; completion requires stated verification.
- `specs/` contains normative, testable capability requirements.

## Working agreement

1. Complete work in dependency order in `tasks.md` and keep evidence current.
2. Adjust the relevant `spec.md` before changing a user-visible behavior.
3. Never use production health data in tests, issues, CI artifacts, or AI prompts beyond the requested scope.
4. The $0 target is a budget guardrail, not a guarantee: free tiers and quotas can change. Configure limits and alerts before accepting real data.
5. This is a personal-use tool, not a medical device. It does not diagnose, prescribe, or provide emergency guidance.

## Status lifecycle

`Draft` -> `Approved` -> `In implementation` -> `Release candidate` ->
`Released`. Release requires all required tasks and specification scenarios to
have automated or recorded manual evidence.
