package plugins

import "github.com/cloudwego/eino/components/tool"

var (
	Plugin = NewPluginRegistry()
)


type PluginRegistry struct {
	AllMcpTool []AgentMcpTool
	FilteredMcpTool []AgentMcpTool
	CommonTool map[string]AgentCommonTool
}

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		AllMcpTool: make([]AgentMcpTool, 0),
		FilteredMcpTool: make([]AgentMcpTool, 0),
		CommonTool: make(map[string]AgentCommonTool),
	}
}

func (r *PluginRegistry) RegisterAllMcpTool(name string, tool AgentMcpTool) {
	r.AllMcpTool = append(r.AllMcpTool, tool)
	r.FilteredMcpTool = append(r.FilteredMcpTool, tool)
}

func (r *PluginRegistry) RegisterFilteredMcpTool(name string, tool AgentMcpTool) {
	r.FilteredMcpTool = append(r.FilteredMcpTool, tool)
}

func (r *PluginRegistry) RegisterCommonTool(name string, tool AgentCommonTool) {
	r.CommonTool[name] = tool
}

func (r *PluginRegistry) GetAllMcpTool() []tool.BaseTool {
	tools := make([]tool.BaseTool, 0, len(r.AllMcpTool))
	for _, tool := range r.AllMcpTool {
		tools = append(tools, tool)
	}
	return tools
}

func (r *PluginRegistry) GetFilteredMcpTool() []tool.BaseTool {
	tools := make([]tool.BaseTool, 0, len(r.FilteredMcpTool))
	for _, tool := range r.FilteredMcpTool {
		tools = append(tools, tool)
	}
	return tools
}

func (r *PluginRegistry) GetCommonTool(name string) []AgentCommonTool {
	tools := make([]AgentCommonTool, 0, len(r.CommonTool))
	for _, tool := range r.CommonTool {
		tools = append(tools, tool)
	}
	return tools
}



