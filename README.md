# ClickUp × HubSpot Integration

Integrates ClickUp tasks with HubSpot using a Go backend and a HubSpot UI extension.

## Overview

- **Go service** — listens for ClickUp webhooks and syncs task data to HubSpot via the HubSpot API
- **HubSpot app** — custom workflow actions that connect HubSpot workflows to ClickUp tasks

## Prerequisites

- Go 1.25+
- A ClickUp API key
- A HubSpot developer account and API token

## Getting Started

1. **Clone the repo**
   ```bash
   git clone https://github.com/prrrprrr/Afosto-Clickup-Hubspot-Integration
   cd your-repo
   ```

2. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Fill in CLICKUP_API_KEY, HUBSPOT_TOKEN, etc.
   ```

3. **Build and run the Go service**
   ```bash
   go build
   ```

4. **Deploy the HubSpot app**
   Follow the [HubSpot CLI docs](https://developers.hubspot.com/docs/platform/workflow-connectors) to upload and install the workflow app in your portal.

## Project Structure

```
.
├── attachments/        # Attachment handling
├── comments/           # Comment handling
├── handlers/           # HTTP request handlers
├── tasks/              # ClickUp task handling
├── tickets/            # HubSpot ticket handling
├── internal/
│   ├── clickup/        # ClickUp client & logic
│   └── hubspot/        # HubSpot client & logic
└── hubspot-app/
    └── app/src/        # HubSpot workflow app source
```

## License

MIT