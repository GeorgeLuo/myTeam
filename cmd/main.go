package main

import (
	"fmt"
	"myTeam/pkg/delegations"
	"myTeam/pkg/messagebuilder"
	"myTeam/pkg/partials"
	"myTeam/pkg/promptbuilder"

	"myTeam/pkg/llmclient/openai"
	"myTeam/pkg/workspace"
)

func main() {

	IDIndex := 1

	agentPromptBuilder := &promptbuilder.AgentPromptBuilderImpl{}

	agentPromptBuilder.SetTopLevelRequirement("You will assist employee 0 in accomplishing his goals as he determines them.")

	agentPromptBuilder.AddOrgMetadata(partials.IdAssignment(IDIndex))
	agentPromptBuilder.AddOrgMetadata(partials.HiringCapabilities(8))
	agentPromptBuilder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/delegation_capabilities.txt"))
	agentPromptBuilder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/defining_responsibilities.txt"))
	agentPromptBuilder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/defining_communication.txt"))

	agentPromptBuilder.AddFunction("Your first responsibility will be of type scheduled. You will provide me a snapshot of the state of your direct reports on a daily basis.")
	agentPromptBuilder.AddFunction("Your second responsibility will be of type message trigger. You will receive a wide range of requests for status of the organization with sometimes granular detail. As you are employee 1, you will have the greatest visibility of all aspects of the organization and capabilities of its parts, and you will be the sole direct report to employee 0.")

	prompt := agentPromptBuilder.ToString()
	description := "A close assistant"
	name := "Corbin"

	authToken := ""

	llmClient := openai.NewOpenAIClient(authToken)
	workspace := workspace.NewWorkspace("default_database")

	assistantID, err := llmClient.CreateAssistant(name, description, prompt)
	if err != nil {
		fmt.Printf("CreateAssistant error: %v\n", err)
		return
	}

	fmt.Printf("assistant id: %v\n", assistantID)
	workspace.AddPersonnel(assistantID, name, description, prompt)

	messageBuilder := &messagebuilder.MessageBuilderImpl{}
	messageBuilder.SetSender("Employee 0")
	messageBuilder.AppendToMessage("I'm starting a project to build an OpenGL visualizer for system processes. Consider this project to be in our portfolio and let's get started.")

	threadID, runID, err := llmClient.SendMessageToAssistantOnNewThread(assistantID, messageBuilder.ToString())
	if err != nil {
		fmt.Printf("SendMessageToAssistantOnNewThread error: %v\n", err)
		return
	}

	fmt.Printf("thread id: %v, run id: %v\n", threadID, runID)

	response, err := llmClient.GetResponse(threadID, runID, 1)
	if err != nil {
		fmt.Printf("GetResponse error: %v\n", err)
		return
	}
	fmt.Printf("message: %v\n", response)

	approvalmbuilder := &messagebuilder.MessageBuilderImpl{}
	approvalmbuilder.SetSender("Employee 0")
	approvalmbuilder.AppendToMessage("These hires make sense.")
	approvalmbuilder.SetResponseParameters("follow the provided documentation to define your response.")
	approvalmbuilder.IncludeTextFromFile("resources/prompt/components/documentation/hiring_api.txt")

	runID, err = llmClient.SendMessageToAssistant(assistantID, threadID, approvalmbuilder.ToString())
	if err != nil {
		fmt.Printf("SendMessageToAssistant error: %v\n", err)
		return
	}

	fmt.Printf("run id: %v\n", runID)

	response, err = llmClient.GetResponse(threadID, runID, 1)
	if err != nil {
		fmt.Printf("GetResponse error: %v\n", err)
		return
	}
	fmt.Printf("message: %v\n", response)

	hiringData, err := delegations.UnmarshalHiringData([]byte(response))
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	for _, role := range hiringData.Roles {

		agentPromptBuilder := &promptbuilder.AgentPromptBuilderImpl{}

		for _, responsibility := range role.Responsibilities {
			agentPromptBuilder.AddFunction(responsibility.Description)
		}
		agentPromptBuilder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/delegation_capabilities.txt"))
		agentPromptBuilder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/defining_responsibilities.txt"))
		agentPromptBuilder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/defining_communication.txt"))

		IDIndex += 1
		agentPromptBuilder.AddOrgMetadata(partials.IdAssignment(IDIndex))
		agentPromptBuilder.AddOrgMetadata(partials.HiringCapabilities(8))

		agentPromptBuilder.SetTopLevelRequirement(role.TopLevelRequirement)

		name := role.Pseudonym
		description := "A direct report to Employee " + fmt.Sprint(hiringData.HiringEmployeeID) + " with title of " + role.Title

		// Attempt to create a new assistant for each role
		assistantID, err := llmClient.CreateAssistant(name, description, agentPromptBuilder.ToString())
		if err != nil {
			fmt.Printf("Error creating assistant for role %s: %v\n", role.Title, err)
			continue
		}

		fmt.Printf("assistant id: %v\n", assistantID)

		// Populate the workspace with personnel records
		// Note: There might be discrepancies in the data required for workspace and what's available in hiringData
		// For instance, 'description' and 'prompt' fields are assumed to be the role's TopLevelRequirement and the first Responsibility's Description respectively
		workspace.AddPersonnel(assistantID, role.Pseudonym, description, agentPromptBuilder.ToString())

		fmt.Printf("Assistant created and added to workspace: %s (ID: %s)\n", role.Title, assistantID)
	}
}
