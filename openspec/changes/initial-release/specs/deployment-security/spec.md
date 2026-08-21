# Deployment and security specification

## Requirement: deployments use short-lived GitHub OIDC credentials

The system SHALL deploy AWS resources from GitHub Actions only by assuming an
AWS role through GitHub OIDC. Trust SHALL limit repository, branch/tag, and
environment claims. Static AWS keys SHALL not be required or stored.

### Scenario: approved production deployment

- **WHEN** an approved job from the configured production environment deploys an allowed revision
- **THEN** it assumes only the scoped deployment role and applies the reviewed production plan

## Requirement: secrets and sensitive data are protected

The system SHALL keep provider tokens, password-verifier material, session keys,
and webhook secrets in managed storage or protected deployment inputs. It SHALL
not commit, expose, or log them.

### Scenario: secret scanning

- **WHEN** a pull request introduces a token-shaped secret outside an approved test mechanism
- **THEN** CI secret scanning fails before merge

## Requirement: controlled-budget observability

The system SHALL emit privacy-preserving availability/error metrics, configure
an alert route and spend threshold, and document the response to cost alerts.

### Scenario: budget threshold reached

- **WHEN** estimated monthly spend reaches the configured threshold
- **THEN** the owner receives an alert and can follow documented steps to investigate and limit spend
