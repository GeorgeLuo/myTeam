package promptbuilder

// AgentPromptBuilder defines the interface for building agent descriptions.
type AgentPromptBuilder interface {
	// SetTopLevelRequirement sets the top-level requirement for the agent.
	SetTopLevelRequirement(text string)
	// AddOrgMetadata adds a factor of localized organization structure.
	AddOrgMetadata(text string)
	// AddFunction adds a specific function responsibility of the agent.
	AddFunction(text string)
	// AddUnderstanding adds a core definition or rule for the agent to know.
	AddUnderstanding(text string)
	// ToString returns the final prompt as a string.
	ToString() string
}
