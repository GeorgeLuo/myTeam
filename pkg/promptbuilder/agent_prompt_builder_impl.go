package promptbuilder

import (
	"myTeam/pkg/partials"
	"strings"
)

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

func (a *AgentPromptBuilderImpl) AddUnderstandingFromFile(filename string) {
	a.understandings = append(a.understandings, partials.LoadFromFile(filename))
}

func (a *AgentPromptBuilderImpl) GetOutline() map[string][]string {
	outline := make(map[string][]string)
	outline["topLevelRequirement"] = []string{a.topLevelRequirement}
	outline["orgMetadata"] = a.orgMetadata
	return outline
}

func (a *AgentPromptBuilderImpl) ToString() string {
	return "[Understandings]\n" + strings.Join(a.understandings, "\n") + "\n\n[Functions]\n" + strings.Join(a.functions, "\n") + "\n\n[Organizational Metadata]\n" + strings.Join(a.orgMetadata, "\n") + "\n\n[Top-Level Requirement]\n" + a.topLevelRequirement
}
