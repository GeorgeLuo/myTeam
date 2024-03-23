package init

import (
	"bufio"
	"fmt"
	"myTeam/pkg/llmclient/openai"
	"myTeam/pkg/promptbuilder"
	"myTeam/pkg/workspace"
	"os"
	"strings"
)

type InitService interface {
	InitWorkspace() error
}

type serviceImpl struct {
	wsDir        string
	openAIApiKey string
}

func NewService(wsDir, openAIApiKey string) InitService {
	return &serviceImpl{
		wsDir:        wsDir,
		openAIApiKey: openAIApiKey,
	}
}

func (s *serviceImpl) InitWorkspace() error {
	fmt.Printf("Initializing a new workspace in directory: %s\n", s.wsDir)
	ws := workspace.NewWorkspace(s.wsDir)
	fmt.Println("New workspace created successfully.")

	fmt.Println("Add default Assistant 1 to workspace? (y/N)")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')

	if strings.TrimSpace(response) == "N" {
		fmt.Println("Custom configuration selected. (Here, you'd add further logic to handle custom configuration.)")
	} else {
		description := "A close assistant"
		name := "Corbin"

		fmt.Println("Default configuration selected. Implementing default Assistant 1 configuration...")
		agentPromptBuilder := promptbuilder.NewAgentPromptBuilderImpl()

		agentPromptBuilder.SetTopLevelRequirement("You will assist employee 0 in accomplishing his goals as he determines them.")

		agentPromptBuilder.AddOrgMetadata("ID", ws.GetNextAssignableID())
		agentPromptBuilder.AddOrgMetadata("MAX_DIRECT_REPORTS", "8")
		agentPromptBuilder.AddOrgMetadata("REPORTING_TO", "0")
		agentPromptBuilder.AddOrgMetadata("NAME", name)

		// delegation
		agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/delegation/delegation_capabilities.txt")
		agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/delegation/defining_responsibilities.txt")
		agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/delegation/defining_communication.txt")
		agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/delegation/multiplier_complement.txt")

		// profile analysis
		agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/profile_analysis/overview.txt")

		// apis
		agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/documentation/courier_api.txt")
		agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/documentation/hiring_api.txt")
		agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/documentation/hiring_api.txt")

		// functions
		agentPromptBuilder.AddFunction("Your first responsibility is to handle a wide range of requests for status of the organization with sometimes granular detail. As you are employee 1, you will have the greatest visibility of all aspects of the organization and capabilities of its parts, and you will be the sole direct report to employee 0.")
		agentPromptBuilder.AddFunction("Your second responsibility is to track and be aware of employee 0's working style vis a vis his personality type. The scope and execution of this analysis is found in provided documentation.")
		agentPromptBuilder.AddFunction("Your third responsibility is to be build the organization to implement projects of the organization.")

		// traits
		agentPromptBuilder.AddTraitFromFile("resources/prompt/components/traits/skeptical_foresight.txt")

		llmClient := openai.NewOpenAIClient(s.openAIApiKey)

		assistantID, err := llmClient.CreateAssistant(name, description, agentPromptBuilder.ToString())
		if err != nil {
			fmt.Printf("CreateAssistant error: %v\n", err)
			return err
		}
		ws.AddPersonnel(ws.GetNextAssignableID(), name, description, agentPromptBuilder, "openAI", map[string]string{"assistant_id": assistantID})
	}
	return nil
}
