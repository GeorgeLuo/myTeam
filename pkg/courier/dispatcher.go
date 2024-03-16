package courier

import (
	"fmt"
	"myTeam/pkg/llmclient"
	"myTeam/pkg/workspace"
	"sync"
)

type DispatcherImpl struct {
	workspace *workspace.Workspace
	llmClient llmclient.LLMClient
	couriers  []*Courier
}

func NewDispatcherImpl(ws *workspace.Workspace, client llmclient.LLMClient) *DispatcherImpl {
	return &DispatcherImpl{
		workspace: ws,
		llmClient: client,
	}
}

type MessagesAndAttachments struct {
	Messages    []Message
	Attachments map[string]string
}

func (d *DispatcherImpl) NewMessage(recipientID string, senderID string, message string) {

	courier := NewCourier(recipientID, d.workspace, d.llmClient)
	courier.AddMessage(senderID, message)

	d.couriers = append(d.couriers, courier)
}

func (d *DispatcherImpl) GetResponsesByRecipient(recipientID string) (senderToMessagesAndAttachments MessagesAndAttachments) {
	senderToMessagesAndAttachments = MessagesAndAttachments{
		Messages:    []Message{},
		Attachments: make(map[string]string),
	}

	for _, courier := range d.couriers {
		attachments, messages := courier.GetMessagesByRecipient(recipientID)

		if len(messages) > 0 || len(attachments) > 0 {
			senderToMessagesAndAttachments.Messages = append(senderToMessagesAndAttachments.Messages, messages...)
			for filename, content := range attachments {
				senderToMessagesAndAttachments.Attachments[filename] = content
			}
		}
	}

	return senderToMessagesAndAttachments
}

func (d *DispatcherImpl) GenerateCouriers() (couriers []Courier) {
	couriers = append(couriers, *NewCourier("recipientID", d.workspace, d.llmClient))
	return couriers
}

func (d *DispatcherImpl) DispatchCouriersAndWait() {
	var wg sync.WaitGroup

	wg.Add(len(d.couriers))

	for _, courier := range d.couriers {
		go func(c *Courier) {
			defer wg.Done()
			_, _, err := c.DispatchAndWait()
			if err != nil {
				fmt.Printf("Failed to dispatch the message: %v\n", err)
			}
		}(courier)
	}

	wg.Wait()
}
