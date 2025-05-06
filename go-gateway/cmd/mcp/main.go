package main

import (
	"github.com/victormacedo996/poc-mcp/internal/config"
	"github.com/victormacedo996/poc-mcp/internal/domain/usecase"
	mcpgo "github.com/victormacedo996/poc-mcp/internal/infrastructure/mcp/mcp-go"
	"github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/router"
)

func main() {
	c := config.GetInstance()
	router := router.NewRouter(c.WebServer)

	llm_interaction_uc := usecase.GetLlmInteractionUsecase()
	mcp := mcpgo.GetMcpgoInstance()

	router.Start(llm_interaction_uc, mcp)

}
