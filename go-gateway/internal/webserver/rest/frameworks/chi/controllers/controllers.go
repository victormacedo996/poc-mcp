package controllers

import (
	"github.com/victormacedo996/poc-mcp/internal/domain/usecase"
	"github.com/victormacedo996/poc-mcp/internal/infrastructure/mcp"
	v1Controller "github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/controllers/v1"
)

type Controllers struct {
	V1Controllers v1Controller.V1Controller
}

func GetControllersInstance(llm_interaction_uc *usecase.LlmInteraction, mcp mcp.Mcp) *Controllers {
	return &Controllers{
		V1Controllers: *v1Controller.GetV1Controller(llm_interaction_uc, mcp),
	}
}
