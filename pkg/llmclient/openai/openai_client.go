package openai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(authToken string) OpenAIClient {
	return OpenAIClient{
		client: openai.NewClient(authToken),
	}
}

func (c *OpenAIClient) CreateAssistant(name string, description string, prompt string) (assistantID string, err error) {
	assistant, err := c.client.CreateAssistant(context.Background(),
		openai.AssistantRequest{
			Name:         &name,
			Description:  &description,
			Model:        openai.GPT4TurboPreview,
			Instructions: &prompt,
		})
	if err != nil {
		fmt.Printf("CreateAssistant error: %v\n", err)
		return
	}

	return assistant.ID, err
}

func (c *OpenAIClient) SendMessageToAssistantOnNewThread(assistantID string, message string) (threadID string, runId string, err error) {
	thread, err := c.client.CreateThread(context.Background(),
		openai.ThreadRequest{
			Messages: []openai.ThreadMessage{
				{
					Role:    "user",
					Content: "I'm starting a project to design a car, consider this project has started and start tracking work.",
				},
			},
		})
	if err != nil {
		fmt.Printf("CreateThread error: %v\n", err)
		return
	}

	runID, err := c.SendMessageToAssistant(assistantID, thread.ID, message)
	if err != nil {
		fmt.Printf("CreateThread error: %v\n", err)
		return
	}
	return thread.ID, runID, err
}

func (c *OpenAIClient) SendMessageToAssistant(assistantID string, threadID string, message string) (runID string, err error) {
	run, err := c.client.CreateRun(context.Background(), threadID,
		openai.RunRequest{
			AssistantID: assistantID,
		})
	if err != nil {
		fmt.Printf("CreateRun error: %v\n", err)
		return
	}

	return run.ID, err
}
