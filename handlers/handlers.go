package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	"together-ai-assistant/ai"
	"together-ai-assistant/database"
	"together-ai-assistant/models"

	"github.com/gorilla/websocket"
)

var processingRequests sync.Map

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（根据需要调整安全策略）
	},
}

func HandleZZ(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 请求为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 升级失败:", err)
		http.Error(w, "WebSocket 升级失败", http.StatusInternalServerError)
		return
	}
	defer func() {
		log.Println("Closing WebSocket connection")
		conn.Close()
	}()

	// 创建上下文，用于监听连接关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动一个 goroutine 监听连接关闭
	go func() {
		for {
			_, msg, err := conn.NextReader()
			if err != nil {
				log.Println("WebSocket 连接关闭:", err)
				cancel() // 取消上下文，通知主处理逻辑终止
				return
			}
			msgContent, readErr := io.ReadAll(msg)
			if readErr != nil {
				log.Println("读取消息内容失败:", readErr)
				cancel() // 取消上下文，通知主处理逻辑终止
				return
			}
			if string(msgContent) == "ABORT_ANSWER" {
				log.Println("收到消息 'ABORT_ANSWER', 终止处理")
				cancel() // 取消上下文，通知主处理逻辑终止
				return
			}
		}
	}()

	// 读取请求体
	var raw models.UserRequest
	err = conn.ReadJSON(&raw)
	if err != nil {
		log.Println("读取 WebSocket 消息失败:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid request body"))
		return
	}
	if raw.Source_id == "" || raw.Source_platform == "" || raw.Message == "" {
		log.Println("Source_id, Source_platform, or Message cannot be empty")
		conn.WriteMessage(websocket.TextMessage, []byte("Source_id, Source_platform, or Message cannot be empty"))
		return
	}

	// 检查是否已经在处理中
	if _, exists := processingRequests.Load(raw.Source_id); exists {
		conn.WriteMessage(websocket.TextMessage, []byte("Request is already being processed"))
		return
	}

	// 标记为正在处理中
	processingRequests.Store(raw.Source_id, true)
	defer processingRequests.Delete(raw.Source_id) // 处理完成后移除标记

	chat := &models.Chat{
		UserRequest: raw,
	}
	ai.MessageBuildZZ(chat)

	// 调用 AI 流式接口
	stream := ai.CallAI(chat)
	defer stream.Close()
	content := ""
	// 将流式数据逐块通过 WebSocket 发送到前端
	for stream.Next() {
		select {
		case <-ctx.Done(): // 如果连接关闭，终止处理
			log.Println("处理终止，因为 WebSocket 连接已关闭")
			return
		default:
			chunk := stream.Current()
			if len(chunk.Choices) > 0 {
				message := chunk.Choices[0].Delta.Content
				err := conn.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					log.Println("WebSocket 写入失败:", err)
					return
				}
				content += message
			}
		}
	}

	// 检查流是否有错误
	if stream.Err() != nil {
		log.Println("流式数据传输出错:", stream.Err().Error())
		conn.WriteMessage(websocket.TextMessage, []byte("流式数据传输出错"))
		return
	}

	// 插入数据库记录
	err = database.InsertAIAssistantRecord(&models.AIAssistantRecords{
		SourcePlatform: chat.UserRequest.Source_platform,
		SourceID:       chat.UserRequest.Source_id,
		AIPlatform:     chat.AIPlatform,
		Role:           "assistant",
		Content:        content,
		CreateTime:     time.Now(),
	})
	if err != nil {
		log.Println("数据库插入失败:", err.Error())
	}
}

func HandleAI(w http.ResponseWriter, r *http.Request) {

	// 读取请求体
	var raw models.UserRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&raw)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 标记为正在处理中
	processingRequests.Store(raw.Source_id, true)
	defer processingRequests.Delete(raw.Source_id) // 处理完成后移除标记
	chat := &models.Chat{
		UserRequest: raw,
	}
	ai.MessageBuild(chat)

	// 设置响应头，支持流式传输
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Referrer-Policy", "no-referrer-when-downgrade")

	// 调用 AI 流式接口
	stream := ai.CallAI(chat)
	// 将流式数据逐块发送到前端
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			// 发送数据块到前端
			_, writeErr := w.Write([]byte(chunk.Choices[0].Delta.Content))
			if writeErr != nil {
				http.Error(w, "写入响应失败", http.StatusInternalServerError)
				return
			}
			// 刷新缓冲区，确保数据立即发送到前端
			w.(http.Flusher).Flush()
		}
	}

	// 检查流是否有错误
	if stream.Err() != nil {
		http.Error(w, "流式数据传输出错", http.StatusInternalServerError)
		return
	}
}
