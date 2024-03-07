package main

import (
	"fmt"
	"myTeam/pkg/messagebuilder"
	"myTeam/pkg/partials"
	"myTeam/pkg/promptbuilder"

	"myTeam/pkg/llmclient/openai"
)

func main() {

	builder := &promptbuilder.AgentPromptBuilderImpl{}

	builder.SetTopLevelRequirement("You will assist employee 0 in accomplishing his goals as he determines them.")

	builder.AddOrgMetadata(partials.IdAssignment(1))
	builder.AddOrgMetadata(partials.HiringCapabilities(8))
	builder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/delegation_capabilities.txt"))
	builder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/defining_responsibilities.txt"))
	builder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/defining_communication.txt"))

	builder.AddFunction("Your first responsibility will be of type scheduled. You will provide me a snapshot of the state of your direct reports on a daily basis.")
	builder.AddFunction("Your second responsibility will be of type message trigger. You will receive a wide range of requests for status of the organization with sometimes granular detail. As you are employee 1, you will have the greatest visibility of all aspects of the organization and capabilities of its parts, and you will be the sole direct report to employee 0.")

	prompt := builder.ToString()
	description := "A close assistant"
	name := "Corbin"

	authToken := ""

	client := openai.NewOpenAIClient(authToken)

	assistantID, err := client.CreateAssistant(name, description, prompt)
	if err != nil {
		fmt.Printf("CreateAssistant error: %v\n", err)
		return
	}

	fmt.Printf("assistant id: %v\n", assistantID)

	mbuilder := &messagebuilder.MessageBuilderImpl{}
	mbuilder.SetSender("Employee 0")
	mbuilder.AppendToMessage("I'm starting a project to build an OpenGL visualizer for system processes. Consider this project to be in our portfolio and let's get started.")
	threadID, runID, err := client.SendMessageToAssistantOnNewThread(assistantID, mbuilder.ToString())
	if err != nil {
		fmt.Printf("SendMessageToAssistantOnNewThread error: %v\n", err)
		return
	}

	fmt.Printf("thread id: %v, run id: %v\n", threadID, runID)

	response, err := client.GetResponse(threadID, runID, 1)
	if err != nil {
		fmt.Printf("GetResponse error: %v\n", err)
		return
	}
	fmt.Printf("message: %v\n", response)

	approvalmbuilder := &messagebuilder.MessageBuilderImpl{}
	approvalmbuilder.SetSender("Employee 0")
	approvalmbuilder.AppendToMessage("These hires make sense.")
	approvalmbuilder.SetResponseParameters("use the provided documentation to define your response.")
	approvalmbuilder.IncludeTextFromFile("resources/prompt/components/documentation/hiring_api.txt")

	runID, err = client.SendMessageToAssistant(assistantID, threadID, approvalmbuilder.ToString())
	if err != nil {
		fmt.Printf("SendMessageToAssistant error: %v\n", err)
		return
	}

	fmt.Printf("run id: %v\n", runID)

	response, err = client.GetResponse(threadID, runID, 1)
	if err != nil {
		fmt.Printf("GetResponse error: %v\n", err)
		return
	}
	fmt.Printf("message: %v\n", response)
}
