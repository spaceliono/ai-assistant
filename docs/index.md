## Together AI Assistant — Documentation

Welcome to the Together AI Assistant documentation. This guide covers public APIs, packages, data models, configuration, and end‑to‑end usage examples.

- Documentation navigation:
  - Endpoints: see [Endpoints](./endpoints.md)
  - Packages: [AI](./ai.md), [Handlers](./handlers.md), [Database](./database.md), [Config](./config.md), [Models](./models.md)

### Quickstart

Prerequisites:
- Go 1.16+
- MySQL 5.7+

Install and run:
```bash
go mod download
go run main.go
```

Configure the service via `config.yml` (see [Config](./config.md)).

### Architecture Overview

- `main.go`: process entrypoint; bootstraps config, database, AI, and HTTP routes
- `handlers/`: HTTP and WebSocket handlers (`/ai`, `/api/front/together_ai_assistant/chat`)
- `ai/`: AI client integration and prompt/message builders
- `database/`: DB initialization and queries (MySQL)
- `models/`: shared types (request types, chat, persistence models)
- `config/`: YAML configuration loader and struct

### Endpoints Summary

- `POST /ai`: streaming HTTP response of assistant content
- `WS /api/front/together_ai_assistant/chat`: WebSocket streaming assistant; send JSON to start; send `ABORT_ANSWER` to cancel

See [Endpoints](./endpoints.md) for full request/response details and examples.

