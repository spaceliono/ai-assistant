# Together AI Assistant

A Go implementation of the Together AI Assistant service.

## Prerequisites

- Go 1.16 or later
- MySQL 5.7 or later

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Configure the database in `config.yml`
4. Run the application:
   ```bash
   go run main.go
   ```

## API Endpoints

- `POST /api/front/together_ai_assistant/chat` - Chat endpoint

## Configuration

Edit `config.yml` to configure:
- Server port
- Database connection details

## License

MIT 