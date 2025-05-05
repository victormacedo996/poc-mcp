package v1Controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/victormacedo996/poc-mcp/internal/config"
	"github.com/victormacedo996/poc-mcp/internal/domain/usecase"
	"github.com/victormacedo996/poc-mcp/internal/infrastructure/llm/ollama"
	"github.com/victormacedo996/poc-mcp/internal/webserver/rest/dto/request"
	jsonResponse "github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/http_response/json"
)

type McpController struct {
	Usecase *usecase.McpClient
}

var mcp_once sync.Once
var mcp_controller *McpController

func GetMcpController(usecase *usecase.McpClient) *McpController {
	if mcp_controller == nil {
		mcp_once.Do(func() {
			mcp_controller = &McpController{
				Usecase: usecase,
			}
		})
	}

	return mcp_controller
}

func (m *McpController) HandleChat(w http.ResponseWriter, r *http.Request) {
	var prompt_request request.Prompt

	err := json.NewDecoder(r.Body).Decode(&prompt_request)
	if err != nil {
		jsonResponse.StatusUnprocessableEntity(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		jsonResponse.StatusInternalServerError(w, r, fmt.Errorf("streaming unsupported"))
		return
	}

	c := config.GetInstance()
	llm_provider := ollama.GetLlmOllamaInstance(c.Ollama)
	out_chan, err_chan := m.Usecase.HandleChat(llm_provider)

	for {
		select {
		case msg, ok := <-out_chan:
			if !ok {
				jsonResponse.StatusNoContent(w, r)
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case err, ok := <-err_chan:
			if !ok {
				return
			}
			jsonResponse.StatusInternalServerError(w, r, err)
			return
		}
	}
}
