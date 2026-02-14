package llm

import (
	"Mengbot/config"
	"Mengbot/internal/core/model"
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func GetRouterPrompt() string {
	return config.Conf.Prompt.RouterPrompt
}

func GetChatPrompt() string {
	return config.Conf.Prompt.ChatPrompt
}

func GetComplexChatPrompt() string {
	return config.Conf.Prompt.ComplexChatPrompt
}

func GetDiaryPrompt() string {
	return config.Conf.Prompt.DiaryPrompt
}

func GetMetadataPrompt() string {
	return config.Conf.Prompt.MetadataPrompt
}

func BuildRouterPrompt(ctx context.Context, m *model.Message) []*schema.Message {
	routerPrompt := GetRouterPrompt()
	template := prompt.FromMessages(schema.GoTemplate,
		schema.SystemMessage(routerPrompt),
		schema.UserMessage("{{.Query}}"),
	)
	messages, err := template.Format(ctx, map[string]any{
		"Query": m.Query,
		"History":   m.History,
	})
	if err != nil {
		return nil
	}
	return messages
}

func BuildEasyChatPrompt(ctx context.Context, m *model.Message) ([]*schema.Message, error) {
	maomengPrompt := GetChatPrompt()
	systemTpl := maomengPrompt
	template := prompt.FromMessages(schema.GoTemplate,
		schema.SystemMessage(systemTpl),
		schema.UserMessage("{{.Query}}"),
	)

	formattedMessages, err := template.Format(ctx, map[string]any{
		"Time":      m.Time,
		"UserName":  m.UserName,
		"UserRole":  m.UserRole,
		"Documents": m.Documents,
		"History":   m.History,
		"Query":     m.Query,
	})

	if err != nil {
		return nil, err
	}
	fmt.Printf("formattedMessages: %v\n", formattedMessages)
	return formattedMessages, nil
}

func BuildComplexChatPrompt(ctx context.Context, m *model.Message) ([]*schema.Message, error) {
	maomengPrompt := GetComplexChatPrompt()
	systemTpl := maomengPrompt
	template := prompt.FromMessages(schema.GoTemplate,
		schema.SystemMessage(systemTpl),
		schema.UserMessage("{{.Query}}"),
	)

	formattedMessages, err := template.Format(ctx, map[string]any{
		"Time":      m.Time,
		"UserName":  m.UserName,
		"UserRole":  m.UserRole,
		"Documents": m.Documents,
		"History":   m.History,
		"Query":     m.Query,
	})

	if err != nil {
		return nil, err
	}
	fmt.Printf("formattedMessages: %v\n", formattedMessages)
	return formattedMessages, nil
}

func BuildDiaryPrompt(ctx context.Context, today string) []*schema.Message {
	// 构建日记提示
	prompt := schema.Message{
		Role:    schema.Assistant,
		Content: config.Conf.Prompt.DiaryPrompt + today,
	}
	return []*schema.Message{&prompt}
}

func BuildMetadataPrompt(ctx context.Context, diaryContent string) []*schema.Message {
	// 构建元数据提示
	template := prompt.FromMessages(schema.GoTemplate,
		schema.SystemMessage(GetMetadataPrompt()),
	)
	formattedMessages, err := template.Format(ctx, map[string]any{
		"diary": diaryContent,
	})
	if err != nil {
		return nil
	}
	return formattedMessages
}
