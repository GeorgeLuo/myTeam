package promptbuilder

import (
	"fmt"
	"myTeam/pkg/partials"
	"strings"
)

type AgentPromptBuilderImpl struct {
	topLevelRequirement string
	orgMetadata         map[string]string
	functions           []string
	understandings      []string
	promptOutline       PromptOutline
}

func NewAgentPromptBuilderImpl() *AgentPromptBuilderImpl {
	return &AgentPromptBuilderImpl{
		orgMetadata: make(map[string]string),
	}
}

func (a *AgentPromptBuilderImpl) SetTopLevelRequirement(text string) {
	a.topLevelRequirement = text

	a.promptOutline.TopLevelRequirement = text
}

func (a *AgentPromptBuilderImpl) AddOrgMetadata(key string, value string) {
	a.orgMetadata[key] = value

	if key == "MAX_DIRECT_REPORTS" {
		a.promptOutline.OrgMetadata.MaxDirectReports = value
	} else if key == "REPORTING_TO" {
		a.promptOutline.OrgMetadata.ReportingTo = append(a.promptOutline.OrgMetadata.ReportingTo, value)
	} else if key == "DIRECT_REPORT" {
		a.promptOutline.OrgMetadata.DirectReports = append(a.promptOutline.OrgMetadata.DirectReports, value)
	} else if key == "ID" {
		a.promptOutline.OrgMetadata.ID = value
	} else if key == "NAME" {
		a.promptOutline.OrgMetadata.Name = value
	}
}

func (a *AgentPromptBuilderImpl) AddFunction(text string) {
	a.functions = append(a.functions, text)
	a.promptOutline.Functions = append(a.promptOutline.Functions, Function{
		Description: text,
	})
}

func (a *AgentPromptBuilderImpl) AddUnderstanding(text string) {
	a.understandings = append(a.understandings, text)
}

func (a *AgentPromptBuilderImpl) AddUnderstandingFromFile(filename string) {
	a.understandings = append(a.understandings, partials.LoadFromFile(filename))
	a.promptOutline.Files = append(a.promptOutline.Files, FileMetadata{
		Category: "CONCEPT",
		Filename: filename,
	})
}

func (a *AgentPromptBuilderImpl) GetOutline() PromptOutline {
	// outline := make(map[string][]string)
	// outline["topLevelRequirement"] = []string{a.topLevelRequirement}
	// outline["orgMetadata"] = mapToStringArray(a.orgMetadata)
	// outline["understandings"] = a.understandings
	return a.promptOutline
}

func mapToStringArray(strStrMap map[string]string) (array []string) {
	for k, v := range strStrMap {
		array = append(array, fmt.Sprintf("%s is %s", k, v))
	}
	return
}

func (a *AgentPromptBuilderImpl) ToString() string {
	return "[Understandings]\n" + strings.Join(a.understandings, "\n") + "\n\n[Functions]\n" + strings.Join(a.functions, "\n") + "\n\n[Organizational Metadata]\n" + strings.Join(mapToStringArray(a.orgMetadata), "\n") + "\n\n[Top-Level Requirement]\n" + a.topLevelRequirement
}
