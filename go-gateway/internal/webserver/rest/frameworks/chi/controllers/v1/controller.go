package v1Controller

import (
	"github.com/victormacedo996/poc-mcp/internal/domain/usecase"
)

type V1Controller struct {
	Health HealthController
	Mcp    LlmInteractionController
}

func GetV1Controller(llm_interaction_uc *usecase.LlmInteraction) *V1Controller {
	return &V1Controller{
		Health: *GetHealthController(),
		Mcp:    *GetLlmInteractionController(llm_interaction_uc),
	}
}
