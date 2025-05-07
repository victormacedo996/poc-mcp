package mcp

import "context"

type Mcp interface {
	FetchTools(ctx context.Context) ([]string, error)
	CallTool(ctx context.Context, tool_name string, tool_arguments map[string]interface{}) ([]string, error)
	Close()
}
