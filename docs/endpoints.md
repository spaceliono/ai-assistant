## API Endpoints

This service exposes two public endpoints: a streaming HTTP endpoint and a WebSocket endpoint.

### POST /ai (HTTP Streaming)

- Method: `POST`
- Path: `/ai`
- Content-Type: `application/json`
- Response: `text/plain`, streamed in chunks

Request body:
```json
{
  "source_platform": "web|cli|...",
  "source_id": "a-user-or-session-id",
  "message": "your prompt"
}
```

Example (curl):
```bash
curl -N -X POST \
  -H 'Content-Type: application/json' \
  -d '{"source_platform":"web","source_id":"user-123","message":"Hello"}' \
  http://localhost:30003/ai
```

Example (browser JavaScript):
```javascript
async function streamAI() {
  const res = await fetch('http://localhost:30003/ai', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      source_platform: 'web',
      source_id: 'user-123',
      message: 'Hello'
    })
  });
  const reader = res.body.getReader();
  const decoder = new TextDecoder();
  for (;;) {
    const { value, done } = await reader.read();
    if (done) break;
    const chunk = decoder.decode(value, { stream: true });
    console.log(chunk);
  }
}
```

Notes:
- CORS is enabled with `Access-Control-Allow-Origin: *`.
- The handler flushes after each chunk; consumers should read incrementally.

### WS /api/front/together_ai_assistant/chat (WebSocket)

- Scheme: `ws://` or `wss://` depending on deployment
- Path: `/api/front/together_ai_assistant/chat`
- Protocol: Text frames

Client flow:
1) Connect WebSocket
2) Send JSON to start processing
3) Receive streamed text messages
4) Optional: send the literal string `ABORT_ANSWER` to cancel the current response

Initial start message (JSON):
```json
{
  "source_platform": "web|cli|...",
  "source_id": "a-user-or-session-id",
  "message": "your prompt"
}
```

Example (browser JavaScript):
```javascript
const ws = new WebSocket('ws://localhost:30003/api/front/together_ai_assistant/chat');

ws.onopen = () => {
  ws.send(JSON.stringify({
    source_platform: 'web',
    source_id: 'user-123',
    message: 'Hi there'
  }));
};

ws.onmessage = (ev) => {
  // Each message is a streamed text chunk
  console.log(ev.data);
};

// To cancel the current answer:
// ws.send('ABORT_ANSWER');
```

Server behavior:
- On first JSON message, the server builds the conversation and starts streaming assistant output.
- If `ABORT_ANSWER` is received, the current request is cancelled.
- After completion, the assistant response is stored in the database.

