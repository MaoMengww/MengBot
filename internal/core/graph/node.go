package graph

import (
	"Mengbot/internal/core/llm"
	"Mengbot/internal/core/prompt"
	"Mengbot/pkg/logger"
	"Mengbot/plugins"
	"context"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

/*
路由节点，用途为分析用户意图，判断是否需要复杂回答
主要用途为减少token消耗
*/
func NewRouterLambda() *compose.Lambda {
	router := compose.InvokableLambda(func(ctx context.Context, input *prompt.Message) (*prompt.Message, error) {
		messages := input.BuildRouterPrompt(ctx)
		resp, err := llm.RouterModel.Generate(ctx, messages)
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
	return router
}

/*
简单回答节点，主要为不需要tool工具
*/
func NewEasyChatLambda() *compose.Lambda {
	easyChatNode := compose.InvokableLambda(func(ctx context.Context, input *prompt.Message) (string, error) {
		messages, err := input.BuildEasyChatPrompt(ctx)
		if err != nil {
			return "", err
		}
		resp, err := llm.ChatModel.Generate(ctx, messages)
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
	return easyChatNode
}

/*
复杂回答节点，主要为需要tool工具
具有react能力
*/

func NewReactAgent(ctx context.Context, tools []tool.BaseTool) (*react.Agent, error) {
	agentTools := compose.ToolsNodeConfig{
		Tools: tools,
		//未知工具调用处理
		UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
			logger.Warnf("LLM Hallucinated tool: %s", name)
			return "（时空钟表空转了两圈，发出咔咔声）喵... 那个时间线好像被封锁了，我没办法直接操作那个功能呢。换个方式试试？", nil
		},
	}
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		Model:       llm.ChatModel,
		ToolsConfig: agentTools,
		MaxStep:     6,
	})
	if err != nil {
		logger.Errorf("创建ReactAgent失败: %v", err)
		return nil, err
	}
	return agent, nil
}

func NewComplexChatLambda(ctx context.Context) *compose.Lambda {
	masterAgent, err := NewReactAgent(ctx, plugins.Plugin.GetFilteredMcpTool())
	UserAgent, err := NewReactAgent(ctx, plugins.Plugin.GetFilteredMcpTool())
	if err != nil {
		logger.Errorf("创建ReactAgent失败: %v", err)
		return nil
	}
	complexChatNode := compose.InvokableLambda(func(ctx context.Context, input *prompt.Message) (string, error) {
		var agent *react.Agent
		if input.UserRole == "master" {
			agent = masterAgent
		} else {
			agent = UserAgent
		}
		messages, err := input.BuildEasyChatPrompt(ctx)
		if err != nil {
			return "", err
		}
		resp, err := agent.Generate(ctx, messages)
		if err != nil {
			return "", err
		}
		logger.Infof("回答用户(complexChat): %s, 输入消耗token:%d(缓存命中:%d), 输出消耗token:%d, 共计消耗token:%d",
			input.UserName,
			resp.ResponseMeta.Usage.PromptTokens,
			resp.ResponseMeta.Usage.PromptTokenDetails.CachedTokens,
			resp.ResponseMeta.Usage.CompletionTokens,
			resp.ResponseMeta.Usage.TotalTokens)
		return resp.Content, nil
	})
	return complexChatNode
}

func NewChatBranch() *compose.GraphBranch {
	branch := compose.NewGraphBranch(func(ctx context.Context, input *prompt.Message) (string, error) {
		switch input.Strategy {
		case "COMPLEX":
			return "complexChat", nil
		case "CHAT":
			return "easyChat", nil
		default:
			return compose.END, nil
		}
	}, map[string]bool{
		"complexChat": true,
		"easyChat":    true,
	})
	return branch
}
