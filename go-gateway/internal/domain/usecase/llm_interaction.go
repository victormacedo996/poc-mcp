package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/victormacedo996/poc-mcp/internal/infrastructure/llm"
)

type LlmInteraction struct {
	TotalPrompt string
	FinalPrompt string
}

func GetLlmInteractionUsecase() *LlmInteraction {
	total_prompt := `
SYSTEM:  
You are an AI assistant with access to the following MCP tools (each tool follows MCP’s JSON-RPC 2.0 schema):

TOOLS:
%s

When a user asks a question, follow these steps:

1. **Analyze** the user’s request.  
2. **Decide**:
- If the question requires external data or actions beyond your training (e.g. database lookup, computation), choose exactly one tool from the above list.
- Otherwise, plan to answer directly without calling any tool.

3. **Output** **strictly** one ARRAY containing JSON OBJECTS (no extra text) in one of these forms:

– **Tool call**  
[
  { "name": <tool_name>, "arguments": { … } }
]

You also have the option to do not call any tool, in which case you should return an empty array without any empty object inside it:

[]

User’s question begins below.

USER: %s
	`

	final_prompt := `
You are an expert assistant.  
You have been given:  

1. The original user question:  
   %s

2. The raw MCP tool outputs (already fetched via MCP’s tools/call), in JSON form:  
   %s

Your task is to consume those MCP results and craft a coherent, self-contained answer to the user.  

Assistant (plan):
1. Parse the user question.
2. Examine the MCP_TOOL_OUTPUTS_JSON; identify relevant entries.
3. Summarize or compare the key fields from those entries.
4. Structure the final response: include context, findings, and recommendations.

Assistant (analysis):
– Refer to each MCP result by its index (e.g. “Result #0 shows …”).  
– Extract the most salient fields (e.g. title, summary, metrics).  
– Note any discrepancies or highlights across results.

Assistant (answer):
Provide a clear, concise answer to “%s that weaves in the MCP results. Cite each result by index.  
	`

	return &LlmInteraction{
		TotalPrompt: total_prompt,
		FinalPrompt: final_prompt,
	}
}

func (m *LlmInteraction) HandleCAsynchat(prompt string, llm_provider llm.LLM) (<-chan string, <-chan error) {
	out_chan := make(chan string)
	err_chan := make(chan error)
	token_count_chan := make(chan int)

	go m.asyncChatLoop(prompt, llm_provider, out_chan, err_chan, token_count_chan)
	go m.countTokens(token_count_chan, err_chan)

	return out_chan, err_chan
}

func (m *LlmInteraction) ChooseToolsToCall(ctx context.Context, mcp_tools []string, prompt string, llm_provider llm.LLM) ([]map[string]interface{}, error) {

	tool_prompt := fmt.Sprintf(m.TotalPrompt, mcp_tools, prompt)

	resp, err := llm_provider.SyncChat(tool_prompt)
	if err != nil {
		return nil, err
	}

	var tools_to_call []map[string]interface{}

	if err := json.Unmarshal([]byte(resp), &tools_to_call); err != nil {
		return nil, err
	}

	return tools_to_call, nil
}

func (m *LlmInteraction) HandleSyncChat(ctx context.Context, user_prompt string, mcp_tools string, llm_provider llm.LLM) (string, error) {
	prompt := fmt.Sprintf(m.FinalPrompt, user_prompt, mcp_tools, user_prompt)

	resp, err := llm_provider.SyncChat(prompt)
	if err != nil {
		return "", err
	}

	return resp, nil

}

func (m *LlmInteraction) asyncChatLoop(prompt string, llmProvider llm.LLM, out_chan chan<- string, err_chan chan<- error, token_count_chan chan<- int) {
	defer close(out_chan)
	defer close(err_chan)
	defer close(token_count_chan)

	var tokenCount int
	out, erro_ch := llmProvider.AsyncChat(prompt)

	for {
		select {
		case msg, ok := <-out:
			if !ok {
				return
			}
			words := len(strings.Fields(msg))
			tokenCount += words
			token_count_chan <- tokenCount
			out_chan <- msg

		case e, ok := <-erro_ch:
			if !ok {
				return
			}
			err_chan <- e
		}
	}
}

func (m *LlmInteraction) countTokens(token_count_chan <-chan int, err_chan <-chan error) {
	var token_count int
	for {
		select {
		case count, ok := <-token_count_chan:
			if !ok {
				fmt.Printf("Total tokens: %d\n", token_count)
				return
			}
			token_count += count
		case err, ok := <-err_chan:
			if !ok {
				fmt.Printf("error: %v", err)
				return
			}
		}
	}
}
