package promptbuilder

import "strings"

type AgentPromptBuilderImpl struct {
	topLevelRequirement string
	orgMetadata         []string
	functions           []string
	understandings      []string
}

func (a *AgentPromptBuilderImpl) SetTopLevelRequirement(text string) {
	a.topLevelRequirement = text
}

func (a *AgentPromptBuilderImpl) AddOrgMetadata(text string) {
	a.orgMetadata = append(a.orgMetadata, text)
}

func (a *AgentPromptBuilderImpl) AddFunction(text string) {
	a.functions = append(a.functions, text)
}

func (a *AgentPromptBuilderImpl) AddUnderstanding(text string) {
	a.understandings = append(a.understandings, text)
}

func (a *AgentPromptBuilderImpl) ToString() string {
	return "[Understandings]\n" + strings.Join(a.understandings, "\n") + "\n\n[Functions]\n" + strings.Join(a.functions, "\n") + "\n\n[Organizational Metadata]\n" + strings.Join(a.orgMetadata, "\n") + "\n\n[Top-Level Requirement]\n" + a.topLevelRequirement
}
