package main

import (
	"fmt"
	"myTeam/pkg/courier"
	openai "myTeam/pkg/llmclient/openai"
	"myTeam/pkg/messagebuilder"
	"myTeam/pkg/partials"
	"myTeam/pkg/promptbuilder"
	"myTeam/pkg/workspace"
)

func main() {

	authToken := ""

	// // interfacing element with LLM to create personnel
	llmClient := openai.NewOpenAIClient(authToken)

	// // workplace tracker
	workspace := workspace.NewWorkspace("default_database")

	// // counter for generating ids
	IDIndex := 1

	// // build the prompt to generate agent 1
	agentPromptBuilder := &promptbuilder.AgentPromptBuilderImpl{}

	agentPromptBuilder.SetTopLevelRequirement("You will assist employee 0 in accomplishing his goals as he determines them.")

	agentPromptBuilder.AddOrgMetadata(partials.IdAssignment(IDIndex))
	agentPromptBuilder.AddOrgMetadata(partials.HiringCapabilities(3))
	agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/delegation_capabilities.txt")
	agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/defining_responsibilities.txt")
	agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/defining_communication.txt")
	agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/documentation/courier_api.txt")
	agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/documentation/hiring_api.txt")

	agentPromptBuilder.AddFunction("Your first responsibility will be of type scheduled. You will provide me a snapshot of the state of your direct reports on a daily basis.")
	agentPromptBuilder.AddFunction("Your second responsibility will be of type message trigger. You will receive a wide range of requests for status of the organization with sometimes granular detail. As you are employee 1, you will have the greatest visibility of all aspects of the organization and capabilities of its parts, and you will be the sole direct report to employee 0.")

	description := "A close assistant"
	name := "Corbin"

	// // create personnel
	assistantID, err := llmClient.CreateAssistant(name, description, agentPromptBuilder.ToString())
	if err != nil {
		fmt.Printf("CreateAssistant error: %v\n", err)
		return
	}

	fmt.Printf("assistant id: %v\n", assistantID)

	// // add to workspace
	workspace.AddPersonnel(fmt.Sprint(IDIndex), name, description, agentPromptBuilder, "openAI", map[string]string{"assistant_id": assistantID})

	// // initial message to begin project
	messageBuilder := &messagebuilder.MessageBuilderImpl{}
	messageBuilder.SetSender("Employee 0")
	messageBuilder.AppendToMessage("I'm starting a project to build an OpenGL visualizer for system processes. Consider this project to be in our portfolio and let's get started. Additionally, please include an attachment to test the courier api.")

	// threadID, runID, err := llmClient.SendMessageToAssistantOnNewThread(assistantID, messageBuilder.ToString())
	// if err != nil {
	// 	fmt.Printf("SendMessageToAssistantOnNewThread error: %v\n", err)
	// 	return
	// }

	// fmt.Printf("thread id: %v, run id: %v\n", threadID, runID)

	// if err := workspace.SetModelMetaDataByID(fmt.Sprint(IDIndex), "thread_id", threadID); err != nil {
	// 	fmt.Println("Failed to set model metadata:", err)
	// 	return
	// }

	courier := courier.NewCourier("1", &workspace, &llmClient)
	courier.AddMessage("0", messageBuilder.ToString())
	_, _, err = courier.DispatchAndWait()
	if err != nil {
		fmt.Printf("Dispatch error: %v\n", err)
		return
	}

	attachments, messages := courier.GetMessagesByRecipient("0")
	for _, message := range messages {
		fmt.Printf("message: %+v\n", message)
	}
	for _, attachment := range attachments {
		fmt.Printf("attachment: %+v\n", attachment)
	}

	// approvalmbuilder := &messagebuilder.MessageBuilderImpl{}
	// approvalmbuilder.SetSender("Employee 0")
	// approvalmbuilder.AppendToMessage("These hires make sense.")
	// approvalmbuilder.SetResponseParameters("follow the provided documentation to define your response.")
	// approvalmbuilder.IncludeTextFromFile("resources/prompt/components/documentation/hiring_api.txt")

	// runID, err = llmClient.SendMessageToAssistant(assistantID, threadID, approvalmbuilder.ToString())
	// if err != nil {
	// 	fmt.Printf("SendMessageToAssistant error: %v\n", err)
	// 	return
	// }

	// fmt.Printf("run id: %v\n", runID)

	// response, err = llmClient.GetResponse(threadID, runID, 1)
	// if err != nil {
	// 	fmt.Printf("GetResponse error: %v\n", err)
	// 	return
	// }
	// fmt.Printf("message: %v\n", response)

	// hiringData, err := delegations.UnmarshalHiringData([]byte(response))
	// if err != nil {
	// 	fmt.Println("Error unmarshalling JSON:", err)
	// 	return
	// }

	// // initialize suggested roles
	// for _, role := range hiringData.Roles {

	// 	agentPromptBuilder := &promptbuilder.AgentPromptBuilderImpl{}

	// 	for _, responsibility := range role.Responsibilities {
	// 		agentPromptBuilder.AddFunction(responsibility.Description)
	// 	}
	// 	agentPromptBuilder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/delegation_capabilities.txt"))
	// 	agentPromptBuilder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/defining_responsibilities.txt"))
	// 	agentPromptBuilder.AddUnderstanding(partials.LoadFromFile("resources/prompt/components/defining_communication.txt"))

	// 	IDIndex += 1
	// 	agentPromptBuilder.AddOrgMetadata(partials.IdAssignment(IDIndex))
	// 	agentPromptBuilder.AddOrgMetadata(partials.HiringCapabilities(8))

	// 	agentPromptBuilder.SetTopLevelRequirement(role.TopLevelRequirement)

	// 	name := role.Pseudonym
	// 	description := "A direct report to Employee " + fmt.Sprint(hiringData.HiringEmployeeID) + " with title of " + role.Title

	// 	// Attempt to create a new assistant for each role
	// 	assistantID, err := llmClient.CreateAssistant(name, description, agentPromptBuilder.ToString())
	// 	if err != nil {
	// 		fmt.Printf("Error creating assistant for role %s: %v\n", role.Title, err)
	// 		continue
	// 	}

	// 	fmt.Printf("assistant id: %v\n", assistantID)

	// 	workspace.AddPersonnel(fmt.Sprint(IDIndex), role.Pseudonym, description, agentPromptBuilder.ToString(), "openAI", map[string]string{"assistant_id": assistantID})
	// 	fmt.Printf("Assistant created and added to workspace: %s (ID: %s)\n", role.Title, assistantID)
	// }

	// confirmbuilder := &messagebuilder.MessageBuilderImpl{}
	// confirmbuilder.SetSender("Employee 0")
	// confirmbuilder.AppendToMessage("Thank you, I have made the hires and they are documented in our workspace database. I've attached the records.")
	// confirmbuilder.AppendToMessage("now that we have more complex communication possibilites, we will be using a courier system as documented in the attached document")
	// confirmbuilder.AppendToMessage("think about the next steps in accomplishing our project, and when the courier makes contact, convey the messages to relevant parties")
	// confirmbuilder.SetResponseParameters("confirm you understand the documentation.")
	// confirmbuilder.IncludeTextFromFile("resources/prompt/components/documentation/courier_api.txt")
	// confirmbuilder.IncludeTextFromFile(workspace.File())

	// runID, err = llmClient.SendMessageToAssistant(assistantID, threadID, confirmbuilder.ToString())
	// if err != nil {
	// 	fmt.Printf("SendMessageToAssistant error: %v\n", err)
	// 	return
	// }

	// response, err = llmClient.GetResponse(threadID, runID, 1)
	// if err != nil {
	// 	fmt.Printf("GetResponse error: %v\n", err)
	// 	return
	// }
	// fmt.Printf("message: %v\n", response)

	// courier := courier.NewCourier("1", &workspace, &llmClient)
	// courier.AddMessage("0", "I'm checking in to test the courier system, did you get this?")
	// threadID, runID, err = courier.Dispatch()
	// if err != nil {
	// 	fmt.Printf("Dispatch error: %v\n", err)
	// 	return
	// }
	// fmt.Printf("courier dispatched on thread: %s with runID: %s\n", threadID, runID)
	// response, err = llmClient.GetResponse(threadID, runID, 1)
	// if err != nil {
	// 	fmt.Printf("GetResponse error: %v\n", err)
	// 	return
	// }
	// fmt.Printf("message: %v\n", response)

	// workspace, err := workspace.LoadFromFile("default_database/workspace.json")
	// if err != nil {
	// 	fmt.Printf("LoadFromFile error: %v\n", err)
	// 	return
	// }
	// fmt.Printf("workspace: %v\n", workspace)

	// courier := courier.NewCourier("1", &workspace, &llmClient)
	// courier.AddMessage("0", "I have some initial requirements which will guide the project. The initial deliverable will code that generates a view of multiple adjacent rectangles of the same width with variable heights that will eventually represent %CPU. This should inform how code should be structured as eventually the model will change but code should be grouped around representing system processes.")

	// response, err := courier.DispatchAndWait()
	// if err != nil {
	// 	fmt.Printf("Dispatch error: %v\n", err)
	// 	return
	// }
	// fmt.Printf("courier returned with: %v\n", response)

	// response, err := llmClient.GetResponse("thread_coEzRBWwfmrM7GimeVMTI8EF", "run_aNWW97lWwJFo2ObwzXOZE4yX", 1)
	// if err != nil {
	// 	fmt.Printf("GetResponse error: %v\n", err)
	// 	return
	// }

	// var dispatchResponse courier.DispatchResponse

	// strData := string(response)
	// if strings.HasPrefix(strData, "```json") && strings.HasSuffix(strData, "```") {
	// 	// Remove the very first and last characters,
	// 	// ensuring removal of the code fences even if the JSON is indented or formatted strangely
	// 	strData = strData[7 : len(strData)-3] // Removes starting ```json and ending ```
	// 	strData = strings.TrimSpace(strData)
	// 	if err := json.Unmarshal([]byte(strData), &dispatchResponse); err != nil {
	// 		fmt.Println("Error unmarshaling response:", err)
	// 		return
	// 	}
	// }

	// for _, message := range dispatchResponse.Messages {
	// 	courier := courier.NewCourier(message.RecipientID, &workspace, &llmClient)

	// 	courier.AddMessage("1", message.Message)

	// 	// Dispatch the message
	// 	response, err := courier.DispatchAndWait()
	// 	if err != nil {
	// 		fmt.Printf("Dispatch error: %v\n", err)
	// 		continue
	// 	}
	// 	fmt.Printf("courier returned with: %v\n", response)
	// }
}
