package llm

import (
	"Mengbot/internal/core/model"
	"Mengbot/pkg/logger"
	"Mengbot/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	//	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/subosito/gotenv"
)

var (
	//	ChatModel *deepseek.ChatModel
	ChatModel      *qwen.ChatModel
	RouterModel    *qwen.ChatModel
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
		APIKey: apiKey,
		Model:  rerankModelName,
	})
	if err != nil {
		logger.Fatal(err)
	}
	EmbeddingModel = embeddingModel
}

func Ptr[T any](v T) *T {
	return &v
}

func CallRouter(ctx context.Context, input *model.Message) (*model.Message, error) {
	messages := BuildRouterPrompt(ctx, input)
	resp, err := RouterModel.Generate(ctx, messages)
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
}

func CallChatModel(ctx context.Context, input *model.Message) (string, error) {
	messages, err := BuildEasyChatPrompt(ctx, input)
	if err != nil {
		return "", err
	}
	resp, err := ChatModel.Generate(ctx, messages)
	if err != nil {
		return "", err
	}
	logger.Infof("回答用户(chat): %s, 输入消耗token:%d(缓存命中:%d), 输出消耗token:%d, 共计消耗token:%d",
		input.UserName,
		resp.ResponseMeta.Usage.PromptTokens,
		resp.ResponseMeta.Usage.PromptTokenDetails.CachedTokens,
		resp.ResponseMeta.Usage.CompletionTokens,
		resp.ResponseMeta.Usage.TotalTokens)
	return resp.Content, nil
}

func CallDiary(ctx context.Context, today string) (string, error) {
	prompt := BuildDiaryPrompt(ctx, today)
	resp, err := ChatModel.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

func CallMetadata(ctx context.Context, diaryContent string) (*model.DiaryMetadata, error) {
	prompt := BuildMetadataPrompt(ctx, diaryContent)
	resp, err := ChatModel.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}
	metadata := &model.DiaryMetadata{}
	if err := json.Unmarshal([]byte(resp.Content), metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func CallDiaryEmbedding(ctx context.Context, meta *model.DiaryMetadata) ([]float32, error) {
	input := fmt.Sprintf("Topic: %s | Keywords: %s | Summary: %s",
		strings.Join(meta.Topic, ", "),
		strings.Join(meta.Keywords, ", "),
		meta.Summary)
	embeddings, err := EmbeddingModel.EmbedStrings(ctx, []string{input})
	if err != nil {
		return nil, err
	}
	// 转换为 float32
	embeddingsF32, err := utils.F64ToF32(embeddings[0])
	if err != nil {
		return nil, err
	}
	return embeddingsF32, nil
}

func CallChatEmbedding(ctx context.Context, chat string) ([]float32, error) {
	embeddings, err := EmbeddingModel.EmbedStrings(ctx, []string{chat})
	if err != nil {
		return nil, err
	}
	// 转换为 float32
	embeddingsF32, err := utils.F64ToF32(embeddings[0])
	if err != nil {
		return nil, err
	}
	return embeddingsF32, nil
}
