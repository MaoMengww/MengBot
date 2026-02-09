package prompt

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type Message struct {
	// 1. 核心交互数据
	Query    string            `json:"query"`    // 用户当前说的话
	History  []*schema.Message `json:"history"`  // 之前的聊天记录
	Strategy string            `json:"strategy"` // 策略 (TOOL / CHAT)

	// 2. 状态与感知 (Context)
	UserName string `json:"user_name"` // 用户名称
	UserRole string `json:"user_role"` // 用户角色 (admin / user)
	Time     string `json:"time"`      // 当前时间

	// 3. RAG / 知识库 (可选)
	Documents []string `json:"documents"` // 检索到的相关文档片段
}

func NewMessage(query string, history []*schema.Message, mood string, userName string, userRole string, time string, documents []string) *Message {
	return &Message{
		Query:     query,
		History:   history,
		UserName:  userName,
		UserRole:  userRole,
		Time:      time,
		Documents: documents,
	}
}

func (m *Message) BuildRouterPrompt() []*schema.Message {
	messages := []*schema.Message{
		schema.SystemMessage(routerPrompt),
		schema.UserMessage(m.Query),
	}
	return messages
}

func (m *Message) BuildEasyChatPrompt(ctx context.Context) ([]*schema.Message, error) {
	systemTpl := fmt.Sprintf(`%s
## 当前状态
- **当前时间：** {{.Time}}
- **对话对象：** {{.UserName}} (身份: {{.UserRole}})
## 知识库回忆 (RAG)
{{if .Documents}}
我回想起了以下相关信息：
{{range .Documents}}- {{.}}
{{end}}
请结合这些信息回答，但不要暴露你在读资料，要像是你自己的记忆。
{{else}}
(无额外记忆检索)
{{end}}
{{if eq .UserRole "admin"}}
(检测到当前用户是主人【梦猫】，请表现得傲娇但亲昵，可以使用更可爱的语气！)
{{else}}
(当前用户是普通路人，请保持礼貌但疏离的猫娘态度，不要过分亲密。)
{{end}}
`, MaomengPrompt)

	template := prompt.FromMessages(schema.FString,
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
	return formattedMessages, nil
}
