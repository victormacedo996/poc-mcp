package llm

type LLM interface {
	AsyncChat(prompt string) (<-chan string, <-chan error)
	SyncChat(prompt string) (string, error)
}
