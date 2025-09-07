## Package `handlers`

HTTP and WebSocket handlers that expose the public API.

### Routes

- `POST /ai` → `HandleAI`
- `WS /api/front/together_ai_assistant/chat` → `HandleZZ`

### Public Functions

```go
func HandleZZ(w http.ResponseWriter, r *http.Request)
```
Upgrades the connection to WebSocket, reads a JSON `models.UserRequest` on open, streams assistant chunks as text messages, and persists the final assistant response. Sending the text `ABORT_ANSWER` cancels processing.

```go
func HandleAI(w http.ResponseWriter, r *http.Request)
```
Accepts a JSON `models.UserRequest` via HTTP POST and streams the assistant response as `text/plain` chunks, flushing after each chunk.

### Usage Examples

- See endpoint examples in [Endpoints](./endpoints.md).

### Security Considerations

- The WebSocket upgrader currently allows all origins (`CheckOrigin: true`). Lock this down for production.
- The server exposes streaming text without authentication; add auth if needed.

