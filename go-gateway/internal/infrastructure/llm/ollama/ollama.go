package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/victormacedo996/poc-mcp/internal/config"
)

type Ollama struct {
	OllamaConfig config.Ollama
}

var ollama_once sync.Once
var ollama *Ollama

func GetLlmOllamaInstance(config config.Ollama) *Ollama {
	if ollama == nil {
		ollama_once.Do(func() {
			ollama = &Ollama{
				OllamaConfig: config,
			}
		})
	}

	return ollama
}

func (o *Ollama) AsyncChat(prompt string) (<-chan string, <-chan error) {

	out_chan := make(chan string)
	err_chan := make(chan error)

	go func() {
		payload := GenerateRequest{
			Model:  o.OllamaConfig.Model,
			Prompt: prompt,
		}
		resp, err := http.Post("http://localhost:11434/api/generate",
			"application/json",
			bytes.NewBuffer([]byte(fmt.Sprintf(`{"model":"%s","prompt":"%s","stream":true}`, payload.Model, payload.Prompt))),
		)
		if err != nil {
			err_chan <- fmt.Errorf("HTTP POST request failed: %v", err)
			return
		}
		decoder := json.NewDecoder(bufio.NewReader(resp.Body))
		defer resp.Body.Close()
		defer close(out_chan)
		defer close(err_chan)
		for {
			var chunk GenerateResponse
			if err := decoder.Decode(&chunk); err != nil {
				if err == io.EOF {
					break
				}
				err_chan <- fmt.Errorf("stream decode error: %v", err)
			}
			if chunk.Token != "" {
				out_chan <- chunk.Token
			}
			if chunk.Done {
				break
			}

			out_chan <- chunk.Response
		}
	}()

	return out_chan, err_chan
}

func (o *Ollama) SyncChat(prompt string) (string, error) {

	fmt.Println(o.OllamaConfig.Model)
	fmt.Println(prompt)
	resp, err := http.Post("http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"model":"%s","prompt":"%s","stream":false}`, o.OllamaConfig.Model, prompt))),
	)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to call llm")
	}

	var generated_response GenerateResponse

	err = json.NewDecoder(bufio.NewReader(resp.Body)).Decode(&generated_response)
	if err != nil {
		return "", err
	}

	return generated_response.Response, nil
}
