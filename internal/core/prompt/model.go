package prompt

import (
	"Mengbot/internal/core/memory"
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type Message struct {
	// 1. 核心交互数据
	Query    string                 `json:"query"`    // 用户当前说的话
	History  []memory.MemoryMessage `json:"history"`  // 之前的聊天记录
	Strategy string                 `json:"strategy"` // 策略 (TOOL / CHAT)

	// 2. 状态与感知 (Context)
	UserName string `json:"user_name"` // 用户名称
	UserRole string `json:"user_role"` // 用户角色 (master / user)
	Time     string `json:"time"`      // 当前时间

	// 3. RAG / 知识库 (可选)
	Documents []string `json:"documents"` // 检索到的相关文档片段
}

func NewMessage(query string, history []memory.MemoryMessage, mood string, userName string, userRole string, time string, documents []string) *Message {
	return &Message{
		Query:     query,
		History:   history,
		UserName:  userName,
		UserRole:  userRole,
		Time:      time,
		Documents: documents,
	}
}

func (m *Message) BuildRouterPrompt(ctx context.Context) []*schema.Message {
	routerPrompt := GetRouterPrompt()
	fmt.Printf("routerPrompt: %s\n", routerPrompt)
	messages := []*schema.Message{
		schema.SystemMessage(routerPrompt),
		schema.UserMessage(m.Query),
	}
	return messages
}

func (m *Message) BuildEasyChatPrompt(ctx context.Context) ([]*schema.Message, error) {
	maomengPrompt := GetChatPrompt()
	fmt.Printf("MaomengPrompt: %s\n", maomengPrompt)
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
