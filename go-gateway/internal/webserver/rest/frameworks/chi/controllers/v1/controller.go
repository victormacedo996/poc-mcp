package v1Controller

import (
	"github.com/victormacedo996/poc-mcp/internal/domain/usecase"
	"github.com/victormacedo996/poc-mcp/internal/infrastructure/mcp"
)

type V1Controller struct {
	Health HealthController
	Mcp    LlmInteractionController
}

func GetV1Controller(llm_interaction_uc *usecase.LlmInteraction, mcp mcp.Mcp) *V1Controller {
	return &V1Controller{
		Health: *GetHealthController(),
		Mcp:    *GetLlmInteractionController(llm_interaction_uc, mcp),
	}
}
