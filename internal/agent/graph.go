package agent

import (
	"Mengbot/internal/agent/prompt"
	"Mengbot/pkg/logger"
	"context"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/compose"
)

var chatModel *deepseek.ChatModel

func InitModel() {
	ctx := context.Background()
	apiKey := os.Getenv("API_KEY")
	baseUrl := os.Getenv("BASE_URL")
	modelName := os.Getenv("MODEL_NAME")
	if apiKey == "" || baseUrl == "" || modelName == "" {
		panic("API_KEY or BASE_URL or MODEL_NAME not set")
	}
	model, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:    apiKey,
		BaseURL:   baseUrl,
		Model: modelName,
	})
	if err != nil {
		logger.Fatal(err)
	}
	chatModel = model
}

func GraphBot(ctx context.Context) {
	g := compose.NewGraph[*prompt.Message, string]()

	router := compose.InvokableLambda(func(ctx context.Context, input *prompt.Message) (*prompt.Message, error) {
		messages := input.BuildRouterPrompt()
		resp, err := chatModel.Generate(ctx, messages)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		logger.Infof("回答用户(router): %s, 输入消耗token:%d(缓存命中:%d), 输出消耗token:%d, 共计消耗token:%d",
			input.UserName,
			resp.ResponseMeta.Usage.PromptTokens,
			resp.ResponseMeta.Usage.PromptTokenDetails.CachedTokens,
			resp.ResponseMeta.Usage.CompletionTokens,
			resp.ResponseMeta.Usage.TotalTokens)
		intent := strings.TrimSpace(strings.ToUpper(resp.Content))
		input.Strategy = intent
		return input, nil
	})

	easyChatNode := compose.InvokableLambda(func(ctx context.Context, input *prompt.Message) (string, error) {
		messages, err := input.BuildEasyChatPrompt(ctx)
		if err != nil {
			return "", err
		}
		resp, err := chatModel.Generate(ctx, messages)
		if err != nil {
			return "", err
		}
		logger.Infof("回答用户(easyChat): %s, 输入消耗token:%d(缓存命中:%d), 输出消耗token:%d, 共计消耗token:%d",
			input.UserName,
			resp.ResponseMeta.Usage.PromptTokens,
			resp.ResponseMeta.Usage.PromptTokenDetails.CachedTokens,
			resp.ResponseMeta.Usage.CompletionTokens,
			resp.ResponseMeta.Usage.TotalTokens)
		return resp.Content, nil
	})


	branch := compose.NewGraphBranch[*prompt.Message](func(ctx context.Context, input *prompt.Message) (string, error) {
		switch input.Strategy {
		case "COMPLEX":
			return "complex_chat", nil
		case "CHAT":
			return "easychat", nil
		default:
			return "", nil
		}
	}, map[string]bool{
		"complex_chat": true,
		"easychat":     true,
	})
}
