package llmclient

type LLMClient interface {
	CreateAssistant(name string, description string, prompt string) (assistantID string, err error)
	SendMessageToAssistantOnNewThread(assistantID, message string) (threadID string, runID string, err error)
	SendMessageToAssistant(assistantID string, threadID string, message string) (runID string, err error)
	GetResponse(threadID string, runID string, limit int) (message string, err error)
	SendMessage(recipientMetadata map[string]string, message string) (threadID string, runID string, err error)
}
