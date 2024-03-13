package courier

import "encoding/json"

type DispatchResponse struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	RecipientID string          `json:"recipient_id"`
	Message     string          `json:"message"`
	Data        json.RawMessage `json:"data"`
	Attachments []Attachment    `json:"attachments"`
}

type Attachment struct {
	Filename  string `json:"filename"`
	UsageType string `json:"filetype"`
}
