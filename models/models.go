package models

import (
	"time"

	"github.com/openai/openai-go"
)

type Chat struct {
	UserRequest UserRequest
	AIPlatform  string
	Messages    []openai.ChatCompletionMessageParamUnion
}

type UserRequest struct {
	Source_platform string `json:"source_platform"`
	Source_id       string `json:"source_id"`
	Message         string `json:"message"`
}

type AIAssistantRecords struct {
	ID             int64     `json:"id" db:"id"`
	SourcePlatform string    `json:"source_platform" db:"source_platform"`
	SourceID       string    `json:"source_id" db:"source_id"`
	AIPlatform     string    `json:"ai_platform" db:"ai_platform"`
	Role           string    `json:"role" db:"role"`
	Content        string    `json:"content" db:"content"`
	CreateTime     time.Time `json:"create_time" db:"create_time"`
}
