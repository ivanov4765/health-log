# Personal deployment specification

## Requirement: GitHub Actions deploys one AWS account

GitHub Actions SHALL deploy Lambda and the single Terraform configuration on a
push to `main` or manual dispatch. It SHALL assume a dedicated AWS role through
GitHub OIDC and SHALL not use stored AWS access keys. The role trust policy
SHALL be limited to this repository and its `main` branch.

### Scenario: main branch deployment

- **WHEN** a Lambda or Terraform change is pushed to `main`
- **THEN** the workflow builds the Linux Lambda artifact, assumes the OIDC role,
  and applies Terraform to the one AWS account

## Requirement: Terraform state is available to CI

The owner SHALL create a private, encrypted, versioned S3 state bucket once
from a local AWS profile. Application Terraform SHALL use that S3 backend with
lockfile locking. The bootstrap state remains ignored locally.

### Scenario: first CI setup

- **WHEN** the owner completes the bootstrap Terraform apply and sets repository
  variables for region, role ARN, and state-bucket name
- **THEN** a workflow run initializes the application Terraform backend without
  creating a separate environment

## Requirement: secrets stay out of source control

The Telegram token, webhook secret, dashboard token, and optional Gemini key
SHALL be named SSM SecureString parameters and excluded from Git. Terraform and
GitHub Actions SHALL receive parameter names only, not secret values.

### Scenario: repository inspection

- **WHEN** the owner searches tracked files for the deployed secret values
- **THEN** none are found

## Requirement: owner can notice unexpected cost

The owner SHALL configure a small AWS Budget or Free Tier alert manually before
using the system with personal records.

### Scenario: cost-alert check

- **WHEN** the owner opens the AWS billing console after setup
- **THEN** the configured alert and email destination are visible
