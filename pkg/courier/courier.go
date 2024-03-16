package courier

import (
	"encoding/json"
	"fmt"
	"myTeam/pkg/llmclient" // Assuming these are the packages where Workspace and LLMClient are defined
	"myTeam/pkg/workspace"
	"strings"
)

type Courier struct {
	recipientID         string
	messages            map[string]string
	workspace           *workspace.Workspace
	llmClient           llmclient.LLMClient
	dispatchResponse    DispatchResponse
	responseAttachments map[string]string
}

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

func (c *Courier) Dispatch() (threadID string, runID string, err error) {
	metadata, exists := c.workspace.Personnel[c.recipientID]
	if !exists {
		fmt.Println("Recipient not found in workspace")
		return threadID, runID, fmt.Errorf("recipient %v not found in workspace", c.recipientID)
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
		return response, attachments, err
	}

	if err = c.workspace.SetModelMetaDataByID(fmt.Sprint(c.recipientID), "thread_id", threadID); err != nil {
		fmt.Println("Failed to set model metadata:", err)
		return
	}

	rawResponse, err := c.llmClient.GetResponse(threadID, runID, 1)
	if err != nil {
		fmt.Printf("GetResponse error: %v\n", err)
		return response, attachments, err
	}

	jsonStart := strings.Index(rawResponse, "<DISPATCH_RESPONSE>") + len("<DISPATCH_RESPONSE>")
	jsonEnd := strings.Index(rawResponse, "</DISPATCH_RESPONSE>")
	if jsonEnd > jsonStart && jsonStart > len("<DISPATCH_RESPONSE>")-1 {
		jsonPart := rawResponse[jsonStart:jsonEnd]
		if err = json.Unmarshal([]byte(jsonPart), &response); err != nil {
			fmt.Println("Error unmarshaling JSON response:", err)
			return response, attachments, err
		}
	} else {
		// TODO: if invalid response format, reply with courier documentation to refresh understanding
		return response, attachments, fmt.Errorf("invalid response format")
	}

	attachments = make(map[string]string)

	attachmentSections := strings.Split(rawResponse[jsonEnd:], "<ATTACHMENT")
	for _, section := range attachmentSections {
		if trimmed := strings.TrimSpace(section); trimmed != "" {
			attachmentRaw := "<ATTACHMENT" + section
			filenameStart := strings.Index(attachmentRaw, "filename=") + len("filename=")
			if filenameStart >= len("filename=") {
				filenameEnd := strings.Index(attachmentRaw[filenameStart:], ">") + filenameStart
				if filenameEnd > filenameStart {
					filename := strings.Trim(attachmentRaw[filenameStart:filenameEnd], "\"")
					contentStart := filenameEnd + 1
					contentEnd := strings.Index(attachmentRaw, "</ATTACHMENT>")
					if contentEnd > contentStart {
						content := attachmentRaw[contentStart:contentEnd]
						attachments[filename] = content
					}
				}
			}
		}
	}

	c.responseAttachments = attachments
	c.dispatchResponse = response
	return response, attachments, nil
}

func (c *Courier) GetMessagesByRecipient(recipientID string) (attachments map[string]string, messages []Message) {
	attachments = make(map[string]string)
	for _, msg := range c.dispatchResponse.Messages {
		if msg.RecipientID == recipientID {
			messages = append(messages, msg)
			for _, attachment := range msg.Attachments {
				if attachmentContent, ok := c.responseAttachments[attachment.Filename]; ok {
					attachments[attachment.Filename] = attachmentContent
				}
			}
		}
	}
	return attachments, messages
}
