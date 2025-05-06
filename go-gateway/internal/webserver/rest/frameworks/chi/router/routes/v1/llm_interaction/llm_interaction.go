package v1LlmInteraction

import (
	"github.com/go-chi/chi/v5"
	v1Controller "github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/controllers/v1"
)

func SetRoutes(r *chi.Mux, controller *v1Controller.LlmInteractionController) {

	r.Route("/chat", func(router chi.Router) {
		router.Post("/", controller.HandleChat)
	})
}
