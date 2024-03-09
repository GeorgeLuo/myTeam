package courier

import (
	"fmt"
	"myTeam/pkg/llmclient" // Assuming these are the packages where Workspace and LLMClient are defined
	"myTeam/pkg/workspace"
)

// Courier struct defined with necessary fields
type Courier struct {
	recipientID string
	messages    map[string]string
	workspace   *workspace.Workspace
	llmClient   llmclient.LLMClient
}

// NewCourier is a constructor function for Courier
func NewCourier(recipientID string, ws *workspace.Workspace, client llmclient.LLMClient) *Courier {
	return &Courier{
		recipientID: recipientID,
		workspace:   ws,
		llmClient:   client,
		messages:    make(map[string]string),
	}
}

func (c *Courier) AddMessage(sender string, message string) {
	c.messages[sender] = message
}

// Dispatch method which sends a message to the specified employee
func (c *Courier) Dispatch() (threadID string, runID string, err error) {
	// Retrieve metadata from workspace for employee
	metadata, exists := c.workspace.Personnel[c.recipientID]
	if !exists {
		fmt.Println("Recipient not found in workspace")
		return
	}

	completeMessage := "This is the Courier\n\n"

	for sender, message := range c.messages {
		completeMessage += "FROM: Employee " + sender + "\n"
		completeMessage += "TO: Employee " + c.recipientID + "\n"
		completeMessage += "MESSAGE: " + message + "\n"
		completeMessage += "\nEND OF MESSAGES\n" + "Please take your time in formulating responses if any, I will receive any messages you have outbound."
	}

	return c.llmClient.SendMessage(metadata.ModelMetadata, completeMessage)
}
