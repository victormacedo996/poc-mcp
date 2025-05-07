package mcpgo

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type Mcpgo struct {
	client *client.Client
}

func NewMcpGo(ctx context.Context, server_base_url string) (Mcpgo, error) {
	ctx2 := context.Background()
	server_url := fmt.Sprintf("%s/sse", server_base_url)
	sse, err := transport.NewSSE(server_url)
	if err != nil {
		return Mcpgo{}, err
	}
	err = sse.Start(ctx2)
	if err != nil {
		return Mcpgo{}, err
	}

	cli := client.NewClient(sse)

	_, err = cli.Initialize(ctx2, mcp.InitializeRequest{})
	if err != nil {
		return Mcpgo{}, err
	}

	return Mcpgo{
		client: cli,
	}, nil
}

func (m Mcpgo) CallTool(ctx context.Context, tool_name string, tool_arguments map[string]interface{}) ([]string, error) {
	call_tool_request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 "json:\"name\""
			Arguments map[string]interface{} "json:\"arguments,omitempty\""
			Meta      *struct {
				ProgressToken mcp.ProgressToken "json:\"progressToken,omitempty\""
			} "json:\"_meta,omitempty\""
		}{
			Name:      tool_name,
			Arguments: tool_arguments,
		},
	}

	resp, err := m.client.CallTool(ctx, call_tool_request)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, content := range resp.Content {
		result = append(result, fmt.Sprintf("%v", content))
	}

	return result, nil
}

func (m Mcpgo) Close() {
	m.client.Close()
}

func (m Mcpgo) FetchTools(ctx context.Context) ([]string, error) {
	_, err := m.client.Initialize(ctx, mcp.InitializeRequest{})
	if err != nil {
		return nil, err
	}
	tools, err := m.client.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, err
	}

	var str_tools []string
	for _, tool := range tools.Tools {
		byte_tool, err := tool.MarshalJSON()
		if err != nil {
			return nil, err
		}

		str_tools = append(str_tools, string(byte_tool))
	}

	return str_tools, nil
}
