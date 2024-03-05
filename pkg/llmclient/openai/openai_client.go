package openai

import (
	"context"
	"fmt"
	"time"

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

	runID, err := c.triggerRun(assistantID, thread.ID)
	if err != nil {
		fmt.Printf("CreateThread error: %v\n", err)
		return
	}
	return thread.ID, runID, err
}

func (c *OpenAIClient) triggerRun(assistantID string, threadID string) (runID string, err error) {
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

func (c *OpenAIClient) SendMessageToAssistant(assistantID string, threadID string, message string) (runID string, err error) {
	_, err = c.client.CreateMessage(context.Background(), threadID,
		openai.MessageRequest{
			Role:    "user",
			Content: message,
		})
	if err != nil {
		fmt.Printf("CreateMessage error: %v\n", err)
		return
	}
	return c.triggerRun(assistantID, threadID)
}

func (c *OpenAIClient) GetResponse(threadID string, runID string, limit int) (message string, err error) {
	for {
		run, err := c.client.RetrieveRun(context.Background(), threadID, runID)
		if err != nil {
			fmt.Printf("RetrieveRun error: %v\n", err)
			// return "", err
		}
		if run.Status == openai.RunStatusCompleted {
			messages, err := c.client.ListMessage(context.Background(), threadID, &limit, nil, nil, nil)
			if err != nil {
				fmt.Printf("ListMessage error: %v\n", err)
				return "", err
			}
			return messages.Messages[0].Content[0].Text.Value, err
		}
		// Sleep for a short duration before checking the status again
		time.Sleep(10 * time.Second)
	}
}
