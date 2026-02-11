package graph

import (
	"Mengbot/internal/core/prompt"
	"Mengbot/pkg/logger"
	"context"

	"github.com/cloudwego/eino/compose"
)

func NewGraph() *compose.Graph[*prompt.Message, string] {
	return compose.NewGraph[*prompt.Message, string]()
}

func GraphBotInit(ctx context.Context) *compose.Runnable[*prompt.Message, string] {
	g := NewGraph()
	// 添加路由节点
	routerLambda := NewRouterLambda()
	err := g.AddLambdaNode("router", routerLambda)
	if err != nil {
		logger.Errorf("添加路由节点失败: %v", err)
		return nil
	}

	// 添加简单回答节点
	easyChatLambda := NewEasyChatLambda()
	err = g.AddLambdaNode("easyChat", easyChatLambda)
	if err != nil {
		logger.Errorf("添加简单回答节点失败: %v", err)
		return nil
	}

	// 添加复杂回答节点
	complexChatLambda := NewComplexChatLambda(ctx)
	err = g.AddLambdaNode("complexChat", complexChatLambda)
	if err != nil {
		logger.Errorf("添加复杂回答节点失败: %v", err)
		return nil
	}

	//添加分支
	branch := NewChatBranch()

	//编排
	err = g.AddEdge(compose.START, "router")
	if err != nil {
		logger.Errorf("添加路由节点到简单回答节点的边失败: %v", err)
		return nil
	}

	err = g.AddBranch("router", branch)
	if err != nil {
		logger.Errorf("添加分支失败: %v", err)
		return nil
	}

	err = g.AddEdge("easyChat", compose.END)
	if err != nil {
		logger.Errorf("添加简单回答节点到结束节点的边失败: %v", err)
		return nil
	}

	err = g.AddEdge("complexChat", compose.END)
	if err != nil {
		logger.Errorf("添加复杂回答节点到结束节点的边失败: %v", err)
		return nil
	}
	r, err := g.Compile(ctx)
	if err != nil {
		logger.Errorf("编译图失败: %v", err)
		return nil
	}
	return &r
}
