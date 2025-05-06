package mcpgo

import (
	"context"
	"fmt"
	"sync"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type Mcpgo struct {
}

var mcpgo_once sync.Once
var mcpgo *Mcpgo

func GetMcpgoInstance() *Mcpgo {
	if mcpgo == nil {
		mcpgo_once.Do(func() {
			mcpgo = &Mcpgo{}
		})
	}

	return mcpgo
}

func (m *Mcpgo) FetchTools(ctx context.Context, server_base_url string) ([]string, error) {
	server_url := fmt.Sprintf("%s/sse", server_base_url)
	sse, err := transport.NewSSE(server_url)
	if err != nil {
		return nil, err
	}
	err = sse.Start(ctx)
	if err != nil {
		return nil, err
	}

	cli := client.NewClient(sse)
	defer cli.Close()

	_, err = cli.Initialize(ctx, mcp.InitializeRequest{})
	if err != nil {
		return nil, err
	}
	tools, err := cli.ListTools(ctx, mcp.ListToolsRequest{})
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
