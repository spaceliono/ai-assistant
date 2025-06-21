package ai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"together-ai-assistant/config"
	"together-ai-assistant/database"
	"together-ai-assistant/models"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/ssestream"
)

var moonshotBaseURL string
var moonshotAPIKey string
var moonshotModel string

var deepseekBaseURL string
var deepseekAPIKey string
var deepseekModel string

func InitAI(cfg *config.Config) {
	moonshotBaseURL = cfg.AI.Moonshot.BaseURL
	moonshotAPIKey = cfg.AI.Moonshot.APIKey
	moonshotModel = cfg.AI.Moonshot.Model
	deepseekBaseURL = cfg.AI.Deepseek.BaseURL
	deepseekAPIKey = cfg.AI.Deepseek.APIKey
	deepseekModel = cfg.AI.Deepseek.Model
}

func MessageBuildZZ(chat *models.Chat) error {
	messages := make([]openai.ChatCompletionMessageParamUnion, 0)
	logs, err := database.GetAIAssistantRecordsBySource(chat.UserRequest.Source_platform, chat.UserRequest.Source_id)
	if err != nil {
		log.Println("数据库查询失败 ", err)
	}
	jobs, err := database.GetJobs()
	if err != nil {
		log.Println("数据库查询失败 ", err)
	}
	if len(logs) > 0 {
		chat.AIPlatform = logs[0].AIPlatform
		for _, log := range logs {
			switch role := log.Role; role {
			case "system":
				messages = append(messages, openai.SystemMessage(log.Content))
			case "user":
				messages = append(messages, openai.UserMessage(log.Content))
			case "assistant":
				messages = append(messages, openai.AssistantMessage(log.Content))
			}
		}
	} else {
		initMessage := `你现在扮演一个共创助手;你负责的内容是帮助用户填写招人完成需求或者设计的表单;
1.你根据与用户的对话帮用户分析需要的人员岗位和市场价格,你要充分的问用户想要什么,给几个选择;
2.合理扩展问题;
3.我们平台有以下职业的注册人员, %s
4.输出内容分两部分,一先输出文本思考过程(不要markdown),二输出json数据,用符号⸻隔开这两部分,不要输出json数据标头,json样例
{
    "jobs":[
    {"jobname":"职位名称(匹配我们平台的职业名称)",
    "jobrequire":"职位要求",
    "jobmarketprice":{
    "pricemode":"报酬模式,0:直接报酬,1:每件商品收益比例分红,2:每件商品收益固定分红,10:每件商品收益比例分红+直接报酬,20:每件商品收益固定分红+直接报酬,21:每件商品收益比例分红+每件商品收益固定分红;210:三种报酬方式都有",
    "directprice":"直接报酬",
    "salerateprice":"每件商品收益比例分红",
    "salefixprice":"每件商品收益固定分红"
    },
    "workers":"建议人数/合作机构数量"
    }],
    "title":"动态地根据对话内容生成合适的标题",
    "desc":"动态地根据对话内容生成合适的描述",
    "topic_tags":["动态地根据对话内容生成合适的话题标签一行一个"]
}`

		initMessage = fmt.Sprintf(initMessage, strings.Join(jobs, ","))
		messages = append(messages, openai.SystemMessage(initMessage))
		chat.AIPlatform = "deepseek"
		err = database.InsertAIAssistantRecord(&models.AIAssistantRecords{
			SourcePlatform: chat.UserRequest.Source_platform,
			SourceID:       chat.UserRequest.Source_id,
			AIPlatform:     chat.AIPlatform,
			Role:           "system",
			Content:        initMessage,
			CreateTime:     time.Now(),
		})
		if err != nil {
			log.Println("数据库插入失败 ", err.Error())
		}
	}
	messages = append(messages, openai.UserMessage(chat.UserRequest.Message))
	chat.Messages = messages
	return nil
}
func MessageBuild(chat *models.Chat) error {
	messages := make([]openai.ChatCompletionMessageParamUnion, 0)
	logs, err := database.GetAIAssistantRecordsBySource(chat.UserRequest.Source_platform, chat.UserRequest.Source_id)
	if err != nil {
		log.Println("数据库查询失败 ", err)
	}

	if len(logs) > 0 {
		chat.AIPlatform = logs[0].AIPlatform
		for _, log := range logs {
			switch role := log.Role; role {
			case "system":
				messages = append(messages, openai.SystemMessage(log.Content))
			case "user":
				messages = append(messages, openai.UserMessage(log.Content))
			case "assistant":
				messages = append(messages, openai.AssistantMessage(log.Content))
			}
		}
	} else {
		initMessage := `你是AI助手，仔细思考，解答问题`
		messages = append(messages, openai.SystemMessage(initMessage))
		chat.AIPlatform = "moonshot"
	}
	messages = append(messages, openai.UserMessage(chat.UserRequest.Message))
	chat.Messages = messages
	return nil
}

func CallAI(chat *models.Chat) (stream *ssestream.Stream[openai.ChatCompletionChunk]) {
	baseUrl := ""
	APIKey := ""
	model := ""
	if chat.AIPlatform == "moonshot" {
		baseUrl = moonshotBaseURL
		APIKey = moonshotAPIKey
		model = moonshotModel
	} else if chat.AIPlatform == "deepseek" {
		baseUrl = deepseekBaseURL
		APIKey = deepseekAPIKey
		model = deepseekModel
	}
	client := openai.NewClient(option.WithBaseURL(baseUrl), option.WithAPIKey(APIKey))
	ctx := context.TODO()
	param := openai.ChatCompletionNewParams{
		Messages:  chat.Messages,
		Seed:      openai.Int(1),
		Model:     model,
		MaxTokens: openai.Int(8 * 1024),
	}

	stream = client.Chat.Completions.NewStreaming(ctx, param)

	err := database.InsertAIAssistantRecord(&models.AIAssistantRecords{
		SourcePlatform: chat.UserRequest.Source_platform,
		SourceID:       chat.UserRequest.Source_id,
		AIPlatform:     chat.AIPlatform,
		Role:           "user",
		Content:        chat.UserRequest.Message,
		CreateTime:     time.Now(),
	})
	if err != nil {
		log.Println("数据库插入失败 ", err.Error())
	}
	return stream
}
