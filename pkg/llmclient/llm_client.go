package llmclient

type LLMClient interface {
	CreateAssistant(name string, description string, prompt string) (assistantID string, err error)
	SendMessageToAssistantOnNewThread(assistant_id string, thread_id string, message string) (assistantID string, err error)
	SendMessageToAssistant(assistantID string, threadID string, message string) (runID string, err error)
}
