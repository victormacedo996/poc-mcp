package controllers

import (
	"github.com/victormacedo996/poc-mcp/internal/domain/usecase"
	v1Controller "github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/controllers/v1"
)

type Controllers struct {
	V1Controllers v1Controller.V1Controller
}

func GetControllersInstance(mcp_uc *usecase.McpClient) *Controllers {
	return &Controllers{
		V1Controllers: *v1Controller.GetV1Controller(mcp_uc),
	}
}
