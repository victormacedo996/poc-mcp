package main

import (
	"github.com/victormacedo996/poc-mcp/internal/config"
	"github.com/victormacedo996/poc-mcp/internal/domain/usecase"
	"github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/router"
)

func main() {
	c := config.GetInstance()
	router := router.NewRouter(c.WebServer)

	mcp_uc := usecase.GetMcpUsecase()
	router.Start(mcp_uc)

}
