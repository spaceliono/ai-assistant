package main

import (
	"log"
	"net/http"
	"os"
	"together-ai-assistant/ai"
	"together-ai-assistant/config"
	"together-ai-assistant/database"
	"together-ai-assistant/handlers"
)

func init() {
	// 捕获标准输出
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("无法创建日志文件:", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
func main() {
	log.Println("程序启动")

	//加载配置
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	ai.InitAI(cfg)

	http.HandleFunc("/api/front/together_ai_assistant/chat", handlers.HandleZZ)
	http.HandleFunc("/ai", handlers.HandleAI)
	err = http.ListenAndServe(":30003", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
