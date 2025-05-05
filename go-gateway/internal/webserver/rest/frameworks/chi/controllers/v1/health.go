package v1Controller

import (
	"net/http"
	"sync"

	dto "github.com/victormacedo996/poc-mcp/internal/webserver/rest/dto/response"
	jsonResponse "github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/http_response/json"
)

type HealthController struct{}

var health_once sync.Once
var health_controller *HealthController

func GetHealthController() *HealthController {
	if health_controller == nil {
		health_once.Do(func() {
			health_controller = &HealthController{}
		})
	}

	return health_controller
}

func (h *HealthController) StatusOk(w http.ResponseWriter, r *http.Request) {
	resp := dto.NewHealthResponse("OK")
	jsonResponse.StatusOk(w, r, resp)

}
func (h *HealthController) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	jsonResponse.StatusMethodNotAllowed(w, r)
}
