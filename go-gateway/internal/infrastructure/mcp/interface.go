package mcp

import "context"

type Mcp interface {
	FetchTools(ctx context.Context, server_base_url string) ([]string, error)
}
