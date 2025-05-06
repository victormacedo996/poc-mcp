package v1Controller

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/victormacedo996/poc-mcp/internal/config"
	"github.com/victormacedo996/poc-mcp/internal/domain/usecase"
	"github.com/victormacedo996/poc-mcp/internal/infrastructure/llm/ollama"
	"github.com/victormacedo996/poc-mcp/internal/infrastructure/mcp"
	jsonResponse "github.com/victormacedo996/poc-mcp/internal/webserver/rest/frameworks/chi/http_response/json"
)

type LlmInteractionController struct {
	Usecase *usecase.LlmInteraction
	Mcp     mcp.Mcp
}

var mcp_once sync.Once
var mcp_controller *LlmInteractionController

func GetLlmInteractionController(usecase *usecase.LlmInteraction, mcp mcp.Mcp) *LlmInteractionController {
	if mcp_controller == nil {
		mcp_once.Do(func() {
			mcp_controller = &LlmInteractionController{
				Usecase: usecase,
				Mcp:     mcp,
			}
		})
	}

	return mcp_controller
}

func (m *LlmInteractionController) HandleChat(w http.ResponseWriter, r *http.Request) {
	// var prompt_request request.Prompt

	// err := json.NewDecoder(r.Body).Decode(&prompt_request)
	// if err != nil {
	// 	fmt.Println("erro")
	// 	jsonResponse.StatusUnprocessableEntity(w, r, err)
	// 	return
	// }

	// w.Header().Set("Content-Type", "text/event-stream")
	// w.Header().Set("Cache-Control", "no-cache")
	// w.Header().Set("Connection", "keep-alive")

	// flusher, ok := w.(http.Flusher)
	// if !ok {
	// 	jsonResponse.StatusInternalServerError(w, r, fmt.Errorf("streaming unsupported"))
	// 	return
	// }

	// c := config.GetInstance()
	// llm_provider := ollama.GetLlmOllamaInstance(c.Ollama)
	// out_chan, err_chan := m.Usecase.HandleCAsynchat("please introduce yourself", llm_provider)

	// for {
	// 	select {
	// 	case msg, ok := <-out_chan:
	// 		if !ok {
	// 			jsonResponse.StatusNoContent(w, r)
	// 			return
	// 		}
	// 		fmt.Fprintf(w, "data: %s\n\n", msg)
	// 		flusher.Flush()
	// 	case err, ok := <-err_chan:
	// 		if !ok {
	// 			return
	// 		}
	// 		jsonResponse.StatusInternalServerError(w, r, err)
	// 		return
	// 	}
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	mcp_tools, err := m.Mcp.FetchTools(ctx, "http://localhost:8000")
	if err != nil {
		jsonResponse.StatusInternalServerError(w, r, err)
		return
	}
	c := config.GetInstance()
	llm_provider := ollama.GetLlmOllamaInstance(c.Ollama)

	resp, err := m.Usecase.HandleSyncChat(ctx, mcp_tools, "How is the wheater in Rio de Janeiro?", llm_provider)
	if err != nil {
		jsonResponse.StatusInternalServerError(w, r, err)
		return
	}

	jsonResponse.StatusOk(w, r, resp)

}
