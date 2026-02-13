package plugins

import "github.com/cloudwego/eino/components/tool"

var (
	Plugin = NewPluginRegistry()
)


type PluginRegistry struct {
	AllMcpTool []tool.BaseTool
	FilteredMcpTool []tool.BaseTool
}

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		AllMcpTool: make([]tool.BaseTool, 0),
		FilteredMcpTool: make([]tool.BaseTool, 0),
	}
}


func (r *PluginRegistry) RegisterAllMcpTool(tool tool.BaseTool) {
	r.AllMcpTool = append(r.AllMcpTool, tool)
	r.FilteredMcpTool = append(r.FilteredMcpTool, tool)
}

func (r *PluginRegistry) RegisterFilteredMcpTool(tool tool.BaseTool) {
	r.FilteredMcpTool = append(r.FilteredMcpTool, tool)
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




