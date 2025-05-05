package router

import (
	"fmt"
	netHttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/victormacedo996/poc-mcp/internal/config"
	"github.com/victormacedo996/poc-mcp/internal/domain/usecase"
	"github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/controllers"
	v1Controller "github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/controllers/v1"
	v1Health "github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/router/routes/v1/health"
	v1Mcp "github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/router/routes/v1/mcp"
)

type Router struct {
	config config.WebServer
	router *chi.Mux
}

func NewRouter(webConfig config.WebServer) *Router {

	return &Router{
		config: webConfig,
		router: chi.NewRouter(),
	}
}

func (r *Router) Start(mcp_uc *usecase.McpClient) {
	ctls := controllers.GetControllersInstance(mcp_uc)
	r.configureRouter()
	r.setUpRoutes(ctls)

	port := fmt.Sprintf(":%v", r.config.SERVER_PORT)
	err := netHttp.ListenAndServe(port, r.router)
	if err != nil {
		panic(err)
	}

}

// func enableCors(r *chi.Mux) {
// 	r.Use(cors.Handler(cors.Options{
// 		AllowedOrigins:   []string{"https://*", "http://*"},
// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
// 		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
// 		ExposedHeaders:   []string{"Link"},
// 		AllowCredentials: false,
// 		MaxAge:           300,
// 	}))
// }

func (r *Router) configureRouter() {
	// enableCors(r.router)
	r.router.Use(middleware.Logger)
	r.router.Use(middleware.Timeout(time.Duration(r.config.TIMEOUT)))
	r.router.Use(middleware.RealIP)
	r.router.Use(middleware.RequestID)
	r.router.Use(middleware.Recoverer)
}

func (r *Router) setUpRoutes(controllers *controllers.Controllers) {
	r.router.Mount("/v1", v1Routes(&controllers.V1Controllers))
}

func v1Routes(v1_controller *v1Controller.V1Controller) *chi.Mux {
	v1 := chi.NewRouter()
	v1Health.SetRoutes(v1, &v1_controller.Health)
	v1Mcp.SetRoutes(v1, &v1_controller.Mcp)

	return v1
}
