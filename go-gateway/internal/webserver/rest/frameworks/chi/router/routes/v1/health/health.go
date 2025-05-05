package v1Health

import (
	"github.com/go-chi/chi/v5"
	v1Controller "github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/controllers/v1"
)

func SetRoutes(r *chi.Mux, controller *v1Controller.HealthController) {

	r.Route("/health", func(router chi.Router) {
		router.Get("/", controller.StatusOk)
		router.Post("/", controller.MethodNotAllowed)
		router.Put("/", controller.MethodNotAllowed)
		router.Delete("/", controller.MethodNotAllowed)
		router.Options("/", controller.MethodNotAllowed)
	})
}
