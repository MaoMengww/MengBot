package plugins

import (
	"Mengbot/config"
	"Mengbot/pkg/logger"
	"context"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/mark3labs/mcp-go/client"
	mmcp "github.com/mark3labs/mcp-go/mcp"
)

func Init() error {
	if err := initMcpTools(config.Conf.MCPPath.Bangumi); err != nil {
		logger.Fatalf("初始化 Bangumi 工具失败: %v", err)
		return err
	}
	return nil
}

func initMcpTools(path string) error {
	// 父级超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := config.Conf.MCPPath.Node
	args := []string{path}

	cli, err := client.NewStdioMCPClient(cmd, nil, args...)
	if err != nil {
		return err
	}

	if err := cli.Start(ctx); err != nil {
		return err
	}

	// Initialize
	initReq := mmcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mmcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mmcp.Implementation{
		Name:    "MengBot",
		Version: "1.0.0",
	}

	initCtx, initCancel := context.WithTimeout(ctx, 8*time.Second)
	defer initCancel()

	_, err = cli.Initialize(initCtx, initReq)
	if err != nil {
		return err
	}
	// 获取工具列表
	einoTools, err := mcp.GetTools(ctx, &mcp.Config{
		Cli: cli,
	})
	if err != nil {
		return err
	}
	logger.Infof("共找到 %d 个 MCP 工具", len(einoTools))

	for _, tool := range einoTools {
		t, err := tool.Info(ctx)
		if err != nil {
			return err
		}
		logger.Infof("注册 MCP 工具: %s, 描述: %s", t.Name, t.Desc)
		Plugin.RegisterAllMcpTool(tool)
		Plugin.RegisterFilteredMcpTool(tool)
	}

	logger.Info("MCP 工具初始化完成")
	return nil
}
