package v1Controller

import "github.com/victormacedo996/poc-mcp/internal/domain/usecase"

type V1Controller struct {
	Health HealthController
	Mcp    McpController
}

func GetV1Controller(mcp_uc *usecase.McpClient) *V1Controller {
	return &V1Controller{
		Health: *GetHealthController(),
		Mcp:    *GetMcpController(mcp_uc),
	}
}
