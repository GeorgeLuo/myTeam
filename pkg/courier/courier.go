package courier

import (
	"encoding/json"
	"fmt"
	"myTeam/pkg/llmclient" // Assuming these are the packages where Workspace and LLMClient are defined
	"myTeam/pkg/workspace"
	"strings"
)

// Courier struct defined with necessary fields
type Courier struct {
	recipientID         string
	messages            map[string]string
	workspace           *workspace.Workspace
	llmClient           llmclient.LLMClient
	dispatchResponse    DispatchResponse
	responseAttachments map[string]string
}

// NewCourier is a constructor function for Courier
func NewCourier(recipientID string, ws *workspace.Workspace, client llmclient.LLMClient) *Courier {
	return &Courier{
		recipientID:         recipientID,
		workspace:           ws,
		llmClient:           client,
		messages:            make(map[string]string),
		dispatchResponse:    DispatchResponse{},
		responseAttachments: make(map[string]string),
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

func (c *Courier) DispatchAndWait() (response DispatchResponse, attachments map[string]string, err error) {
	threadID, runID, err := c.Dispatch()
	if err != nil {
		fmt.Printf("Dispatch error: %v\n", err)
		return
	}
	fmt.Printf("courier dispatched on thread: %s with runID: %s\n", threadID, runID)

	if err := c.workspace.SetModelMetaDataByID(fmt.Sprint(c.recipientID), "thread_id", threadID); err != nil {
		fmt.Println("Failed to set model metadata:", err)
	}

	rawResponse, err := c.llmClient.GetResponse(threadID, runID, 1)
	if err != nil {
		fmt.Printf("GetResponse error: %v\n", err)
		return
	}
	// Extract JSON part
	jsonStart := strings.Index(rawResponse, "<DISPATCH_RESPONSE>") + len("<DISPATCH_RESPONSE>")
	jsonEnd := strings.Index(rawResponse, "</DISPATCH_RESPONSE>")
	if jsonEnd > jsonStart && jsonStart > len("<DISPATCH_RESPONSE>")-1 {
		jsonPart := rawResponse[jsonStart:jsonEnd]
		if err = json.Unmarshal([]byte(jsonPart), &response); err != nil {
			fmt.Println("Error unmarshaling JSON response:", err)
			return
		}
	} else {
		fmt.Println("Invalid response format")
		return
	}
	// Initialize attachments map
	attachments = make(map[string]string)
	// Extract attachments
	attachmentStart := jsonEnd + len("</DISPATCH_RESPONSE>")
	attachmentSection := rawResponse[attachmentStart:]
	if strings.Contains(attachmentSection, "<ATTACHMENT") {
		// Assuming only one attachment for simplicity. Multiple attachments handling might require a loop.
		filenameStart := strings.Index(attachmentSection, "filename=") + len("filename=") - 1
		filenameEnd := strings.Index(attachmentSection[filenameStart:], ">") + filenameStart + 1
		filename := attachmentSection[filenameStart:filenameEnd]
		filename = filename[1 : len(filename)-1] // Remove surrounding quotes
		attachmentContentStart := filenameEnd + 1
		attachmentContentEnd := strings.Index(attachmentSection, "</ATTACHMENT>")
		attachmentContent := attachmentSection[attachmentContentStart:attachmentContentEnd]
		attachments[filename] = attachmentContent
	}
	c.responseAttachments = attachments
	c.dispatchResponse = response
	return
}

func (c *Courier) GetMessagesByRecipient(recipientID string) (attachments []string, messages []Message) {
	for _, msg := range c.dispatchResponse.Messages {
		if msg.RecipientID == recipientID {
			messages = append(messages, msg)
			for _, attachment := range msg.Attachments {
				attachments = append(attachments, c.responseAttachments[attachment.Filename])
			}
		}
	}
	return attachments, messages
}
