package messagebuilder

import (
	"myTeam/pkg/partials"
	"strings"
)

type MessageBuilderImpl struct {
	sender             string
	messageParts       []string
	documents          []string
	responseParameters string
}

func (a *MessageBuilderImpl) SetSender(text string) {
	a.sender = text
}

func (a *MessageBuilderImpl) AppendToMessage(text string) {
	a.messageParts = append(a.messageParts, text)
}

func (a *MessageBuilderImpl) IncludeTextFromFile(filename string) {
	a.documents = append(a.documents, partials.LoadFromFile(filename))
}

func (a *MessageBuilderImpl) SetResponseParameters(text string) {
	a.responseParameters = text
}

func (m *MessageBuilderImpl) ToString() string {
	complete := ""
	if len(m.sender) > 0 {
		complete += "hello, this is " + m.sender + "\n\n"
	}
	complete += strings.Join(m.messageParts, "\n")
	if len(m.responseParameters) > 0 {
		complete += "\n\nProduce your response with the following expectations:\n\n"
		complete += m.responseParameters
	}
	if len(m.documents) > 0 {
		complete += "\n\nThe following is documentation that will be helpful as reference:\n\n"
		complete += strings.Join(m.documents, "\n")
	}
	return complete
}
