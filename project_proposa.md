OpenSpec Proposal: Project health-log
Version: 1.0.0

Project Name: health-log

Architecture: Event-Driven Serverless / Jamstack

Target Platform: AWS (Lambda, DynamoDB), GitHub (Pages, Actions), Telegram API, Google Gemini API

1. Executive Summary
health-log is an open-source, serverless health and medication tracking platform designed to offer zero-maintenance, $0-cost personal health logging with conversational AI analysis.

The system provides a two-way interface for personal health management:

Low-Friction Capture: A Telegram Bot webhook for real-time, on-the-go logging of symptom severity, medication intake, timestamps, and notes.

Analysis & Dashboard: A static, single-page web dashboard hosted on GitHub Pages that visualizes historical trends and provides an interactive interface to query health data using the Google Gemini API.

2. Target Persona & Core Capabilities
Target User
Individuals seeking a privacy-focused, zero-recurring-cost tool to log recurring health events (e.g., headaches, daily medication adherence, blood pressure) without relying on closed third-party SaaS subscriptions.

Core Capabilities
Interactive Capture: Instant event creation directly inside Telegram using custom commands or natural language text parsed via AWS Lambda.

Proactive Reminders: Automated notification dispatch via AWS EventBridge to remind users to log daily doses or check in on ongoing symptoms.

Contextual AI Intelligence: Retrieval-Augmented Generation (RAG) light setup where personal DynamoDB logs are injected into the Google Gemini API prompt to identify health patterns, medication frequency, or symptom correlates.

Secure Web Dashboard: A lightweight client-side application featuring interactive visualization (e.g., Chart.js) and token- or password-protected access to backend APIs.

Zero-Touch Infrastructure: Infrastructure-as-Code (IaC) deployment via Terraform and automated CI/CD via GitHub Actions using keyless OpenID Connect (OIDC) authentication.

3. System Architecture & Components
                   [ Telegram App ]
                          │
                          ▼ (Webhook / HTTP POST)
[ GitHub Pages ] ──► [ AWS Lambda ] ──► [ Amazon DynamoDB ]
   (Frontend)      (Go / al2023)  ◄──      (Table)
                          │
                          ▼ (HTTP POST)
                 [ Google Gemini API ]
Component Breakdown
A. Ingestion Layer (Telegram Bot)
Webhook Endpoint: AWS Lambda Function URL handling incoming JSON updates from the Telegram Bot API.

Payload Handling: Parses user commands (/log, /stats, /help) and structured key-value messages into structured schema attributes.

B. Storage Layer (Amazon DynamoDB)
Table Name: HeadacheLogs (or HealthLogs)

Billing Mode: Pay-Per-Request (PAY_PER_REQUEST).

Partition Key (PK): UserId (String) - Enables multi-user isolation.

Sort Key (SK): Timestamp (String, ISO 8601 UTC) - Enables fast range queries over date/time intervals.

C. Compute Layer (AWS Lambda)
Runtime: Custom Runtime on Amazon Linux 2023 (provided.al2023).

Binary: Compiled Go (bootstrap) binary optimized for cold-start performance (<100ms execution times).

Environment Variables:

TELEGRAM_BOT_TOKEN: Secret token provided by @BotFather.

GEMINI_API_KEY: API authentication key for Google AI Studio.

DASHBOARD_SECRET: Secret hash/passcode for HTTP basic authorization on web dashboard endpoints.

D. Intelligence Layer (Google Gemini API)
Model: gemini-2.5-flash

Workflow: When requested by the dashboard or Telegram interface, Lambda fetches raw log arrays from DynamoDB, constructs a system prompt containing the structured JSON array as context, attaches the user query, and proxies the query to the Gemini REST endpoint.

E. Presentation Layer (GitHub Pages)
Stack: Plain HTML5, CSS3, ES6+ JavaScript (no heavy frameworks required).

Distribution: Hosted on GitHub Pages (https://<username>.github.io/health-log/).

Communication: Makes asynchronous fetch() calls directly to the AWS Lambda Function URL using CORS-enabled headers.

4. Technical Constraints & Non-Functional Requirements
Financial Constraint
Operating Cost: Must strictly remain $0.00/month by leveraging standard AWS Free Tier resources, GitHub Pages hosting, and free API tiers from Google AI Studio and Telegram.

Performance & Operational Requirements
Cold Start Latency: Go binary execution latency must remain under 500ms to allow smooth user interaction within Telegram.

Data Security: No raw API keys (AWS, Gemini, or Telegram) stored in client-side code, git tracking, or raw commit history.

Deployment Security: Deployment to AWS must rely exclusively on short-lived OIDC IAM tokens rather than static AWS Access Keys stored in GitHub Repository Secrets.
