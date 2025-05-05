package usecase

import (
	"github.com/victormacedo996/poc-mcp/internal/infrastructure/llm"
)

type McpClient struct {
}

func GetMcpUsecase() *McpClient {
	return &McpClient{}
}

func (m *McpClient) HandleChat(llm_provider llm.LLM) (<-chan string, <-chan error) {
	out_chan := make(chan string)
	err_chan := make(chan error)

	go func() {
		defer close(out_chan)
		defer close(err_chan)

		out, err := llm_provider.AsyncChat("Hello, how are you?")

		for {
			select {
			case msg, ok := <-out:
				if !ok {
					return
				}
				out_chan <- msg
			case err, ok := <-err:
				if !ok {
					return
				}
				err_chan <- err
			}
		}
	}()

	return out_chan, err_chan
}
