# Dashboard access specification

## Requirement: dashboard data needs the personal access token

Public GitHub Pages files SHALL not contain health data or the access token.
Every event, export, or insight request SHALL include
`X-Health-Log-Token`, which Lambda compares with the SSM-stored token.

### Scenario: visitor opens the public dashboard

- **WHEN** someone loads the GitHub Pages URL without entering the token
- **THEN** they can see only the empty sign-in prompt and no health data

### Scenario: invalid token request

- **WHEN** a request has no token or a wrong token
- **THEN** the API rejects it and returns no event data

## Requirement: dashboard offers a minimal personal view

The dashboard SHALL let the owner enter the token, view the last 30 days of
events, view a simple chart, and export a selected range. It does not edit or
delete events in v0.1.0.

### Scenario: owner views history

- **WHEN** the owner enters the correct token
- **THEN** the dashboard loads the recent list and chart from the API
