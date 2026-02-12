package llm

import (
	"Mengbot/pkg/logger"
	"context"
	"os"
	"time"

	//	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/subosito/gotenv"
)

var (
	//	ChatModel *deepseek.ChatModel
	ChatModel   *qwen.ChatModel
	RouterModel *qwen.ChatModel
	EmbeddingModel *dashscope.Embedder

)

func InitModel() {
	// Init ChatModel
	ctx := context.Background()
	if err := gotenv.Load(); err != nil {
		logger.Fatal(err)
	}


	apiKey := os.Getenv("API_KEY")
	baseUrl := os.Getenv("BASE_URL")
	modelName := os.Getenv("MODEL_NAME")
	routerModelName := os.Getenv("ROUTER_MODEL_NAME")
	rerankModelName := os.Getenv("RERANK_MODEL_NAME")


	if apiKey == "" || baseUrl == "" || modelName == "" || routerModelName == "" {
		panic("API_KEY or BASE_URL or MODEL_NAME or ROUTER_MODEL_NAME not set")
	}

	//init chat model
	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:      apiKey,
		BaseURL:     baseUrl,
		Model:       modelName,
		Temperature: Ptr(float32(0.8)),
		Timeout:     60 * time.Second,
	})
	if err != nil {
		logger.Fatal(err)
	}
	ChatModel = chatModel

	//init router model
	routerModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:      apiKey,
		BaseURL:     baseUrl,
		Model:       routerModelName,
		MaxTokens:   Ptr(50),
		Temperature: Ptr(float32(0.0)),
		Timeout:     5 * time.Second,
	})
	if err != nil {
		logger.Fatal(err)
	}
	RouterModel = routerModel

	//init embedding model
	embeddingModel, err := dashscope.NewEmbedder(ctx, &dashscope.EmbeddingConfig{
		APIKey:  apiKey,
		Model:   rerankModelName,
	})
	if err != nil {
		logger.Fatal(err)
	}
	EmbeddingModel = embeddingModel
}

func Ptr[T any](v T) *T {
	return &v
}
