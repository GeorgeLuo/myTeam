package messagebuilder

// MessageBuilder defines the interface for building messages
type MessageBuilder interface {
	// SetSender to place the identity of sender at beginning of message
	SetSender(text string)
	// AppendToMessage to place text at end of message in progress
	AppendToMessage(text string)
	// IncludeTextFromFile to copy and paste text from file into message
	IncludeTextFromFile(filename string)
	// SetResponseParameters to hint at format of response
	SetResponseParameters(text string)
}
