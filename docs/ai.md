## Package `ai`

AI client integration for Moonshot and Deepseek, prompt/message building, and streaming completions.

### Public Functions

```go
func InitAI(cfg *config.Config)
```
Initializes AI provider settings from config (base URLs, API keys, models).

```go
func MessageBuildZZ(chat *models.Chat) error
```
Builds a conversation for the WebSocket workflow. If prior logs exist for the same `source_platform` and `source_id`, they are replayed. Otherwise, a system prompt is created and persisted. Appends the user's message.

```go
func MessageBuild(chat *models.Chat) error
```
Builds a conversation for the HTTP streaming workflow. Replays past conversation if present, otherwise seeds a minimal system prompt.

```go
func CallAI(chat *models.Chat) *ssestream.Stream[openai.ChatCompletionChunk]
```
Creates a streaming completion using the configured provider and returns a stream for incremental consumption.

### Data Flow

1. Initialize from `config.Config` via `InitAI`
2. Build messages with `MessageBuildZZ` or `MessageBuild`
3. Call `CallAI` and iterate the returned stream, handling chunks until completion or error

### Example: Streaming a Completion

```go
package main

import (
    "fmt"
    "together-ai-assistant/ai"
    "together-ai-assistant/config"
    "together-ai-assistant/models"
)

func main() {
    cfg, _ := config.LoadConfig("config.yml")
    ai.InitAI(cfg)

    chat := &models.Chat{
        UserRequest: models.UserRequest{
            Source_platform: "web",
            Source_id:       "example-1",
            Message:         "Tell me a joke",
        },
    }
    _ = ai.MessageBuild(chat)
    stream := ai.CallAI(chat)
    defer stream.Close()

    for stream.Next() {
        if len(stream.Current().Choices) > 0 {
            fmt.Print(stream.Current().Choices[0].Delta.Content)
        }
    }
    if stream.Err() != nil {
        fmt.Println("stream error:", stream.Err())
    }
}
```

### Notes

- Provider selection (Moonshot vs Deepseek) is based on prior conversation records or defaulted on first request.
- `CallAI` also persists the user message to the database.

