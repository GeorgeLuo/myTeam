package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"myTeam/pkg/courier"
	"myTeam/pkg/delegations"
	"myTeam/pkg/llmclient/openai"
	"myTeam/pkg/messagebuilder"
	"myTeam/pkg/promptbuilder"
	"myTeam/pkg/workspace"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <command> [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "chat":
		// Define a new flag set for the chat command
		chatCmd := flag.NewFlagSet("chat", flag.ExitOnError)
		wsFilePath := chatCmd.String("wsfile", "", "Path to the workspace JSON file")
		senderID := chatCmd.String("as", "", "The sender ID")
		recipientID := chatCmd.String("with", "", "The recipient ID")
		openAIApiKey := chatCmd.String("openaiapikey", "", "Your openai api key")

		// Parse flags specific to the chat command
		if err := chatCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Error parsing flags for chat command: %v\n", err)
			os.Exit(1)
		}

		if *wsFilePath == "" || *senderID == "" || *recipientID == "" || *openAIApiKey == "" {
			fmt.Println("Usage: main.go chat -with <recipientID> -as <senderID> -wsfile <workspace.json> -openaiapikey <api key>")
			os.Exit(1)
		}

		// Load the workspace from the provided JSON file
		ws, err := workspace.LoadFromFile(*wsFilePath)
		if err != nil {
			fmt.Printf("Failed to load workspace: %v\n", err)
			return
		}

		llmClient := openai.NewOpenAIClient(*openAIApiKey)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("You are now chatting with Employee %s as %s. Type your message below:\n", *recipientID, *senderID)

		for {
			fmt.Printf("Message (Employee %s): ", *senderID)
			scanned := scanner.Scan()
			if !scanned {
				if err := scanner.Err(); err != nil {
					fmt.Printf("Error reading from input: %v\n", err)
				}
				break
			}
			input := scanner.Text()

			// Check if the user wants to quit
			if input == "/quit" {
				fmt.Println("Exiting chat...")
				break
			}

			courier := courier.NewCourier(*recipientID, &ws, &llmClient)

			// Build the message to be dispatched
			messageBuilder := &messagebuilder.MessageBuilderImpl{}
			messageBuilder.AppendToMessage(input)

			courier.AddMessage(*senderID, messageBuilder.ToString())

			_, _, err := courier.DispatchAndWait()
			if err != nil {
				fmt.Printf("Failed to dispatch the message: %v\n", err)
				continue
			}

			attachments, messages := courier.GetMessagesByRecipient(*senderID)
			for _, message := range messages {
				fmt.Printf("%+v\n", message.Message)
				if len(attachments) > 0 {

					reviewingAttachments := true

					for reviewingAttachments {
						fmt.Printf("*** SYSTEM ***\nfound %d attachment(s)\n", len(attachments))
						for _, attachment := range message.Attachments {
							fmt.Printf("%s\n", attachment.Filename)
						}
						fmt.Println("*** SYSTEM ***")
						scanned := scanner.Scan()
						if !scanned {
							if err := scanner.Err(); err != nil {
								fmt.Printf("Error reading from input: %v\n", err)
							}
							break
						}
						input := scanner.Text()

						// Check if the user pressed <ENTER> to exit attachment review
						if input == "" {
							fmt.Println("Exiting attachment review...")
							break
						}

						fmt.Printf("%s", attachments[input])
					}
				}
				if message.DataSchemaType == "HIRING_RECOMMENDATIONS" {
					fmt.Println("*** SYSTEM ***\nfound agent recommendations in message")
					var hiringData delegations.HiringData
					err := json.Unmarshal(message.Data, &hiringData)
					if err != nil {
						fmt.Printf("Failed to unmarshal hiring recommendations: %v\n", err)
						continue
					}
					hiringDataJSON, _ := json.MarshalIndent(hiringData, "", "  ")
					fmt.Printf("%s\n\n", hiringDataJSON)
					fmt.Println("Approve recommendations to workspace? (y/N)")
					fmt.Println("*** SYSTEM ***")

					reader := bufio.NewReader(os.Stdin)
					approveResponse, _ := reader.ReadString('\n')

					if strings.TrimSpace(approveResponse) == "y" {
						for _, role := range hiringData.Roles {

							description := "A direct report to Employee " + fmt.Sprint(hiringData.HiringEmployeeID) + " with title of " + role.Title
							name := role.Pseudonym

							agentPromptBuilder := promptbuilder.NewAgentPromptBuilderImpl()

							agentPromptBuilder.SetTopLevelRequirement(role.TopLevelRequirement)

							agentPromptBuilder.AddOrgMetadata("ID", ws.GetNextAssignableID())
							for _, reportsTo := range role.ReportsTo {
								agentPromptBuilder.AddOrgMetadata("REPORTING_TO", reportsTo)
							}
							agentPromptBuilder.AddOrgMetadata("NAME", name)

							agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/delegation_capabilities.txt")
							agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/defining_responsibilities.txt")
							agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/defining_communication.txt")
							agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/documentation/courier_api.txt")
							agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/documentation/hiring_api.txt")

							for _, responsibility := range role.Responsibilities {
								agentPromptBuilder.AddFunction(responsibility.Description)
							}

							assistantID, err := llmClient.CreateAssistant(name, description, agentPromptBuilder.ToString())
							if err != nil {
								fmt.Printf("CreateAssistant error: %v\n", err)
								return
							}
							ws.AddPersonnel(ws.GetNextAssignableID(), name, description, agentPromptBuilder, "openAI", map[string]string{"assistant_id": assistantID})
						}
					}
				}
			}
		}

	case "init":
		initCmd := flag.NewFlagSet("init", flag.ExitOnError)
		wsDir := initCmd.String("wsdir", "", "The directory where the new workspace will be created")
		openAIApiKey := initCmd.String("openaiapikey", "", "Your openai api key")

		if err := initCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Error parsing flags for init command: %v\n", err)
			os.Exit(1)
		}

		if *wsDir == "" || *openAIApiKey == "" {
			fmt.Println("Usage: main.go init -wsdir <directory> -openaiapikey <api key>")
			os.Exit(1)
		}

		fmt.Printf("Initializing a new workspace in directory: %s\n", *wsDir)
		ws := workspace.NewWorkspace(*wsDir)
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

			agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/delegation_capabilities.txt")
			agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/defining_responsibilities.txt")
			agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/defining_communication.txt")
			agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/documentation/courier_api.txt")
			agentPromptBuilder.AddUnderstandingFromFile("resources/prompt/components/documentation/hiring_api.txt")

			agentPromptBuilder.AddFunction("Your first responsibility will be of type scheduled. You will provide me a snapshot of the state of your direct reports on a daily basis.")
			agentPromptBuilder.AddFunction("Your second responsibility will be of type message trigger. You will receive a wide range of requests for status of the organization with sometimes granular detail. As you are employee 1, you will have the greatest visibility of all aspects of the organization and capabilities of its parts, and you will be the sole direct report to employee 0.")

			llmClient := openai.NewOpenAIClient(*openAIApiKey)

			assistantID, err := llmClient.CreateAssistant(name, description, agentPromptBuilder.ToString())
			if err != nil {
				fmt.Printf("CreateAssistant error: %v\n", err)
				return
			}
			ws.AddPersonnel(ws.GetNextAssignableID(), name, description, agentPromptBuilder, "openAI", map[string]string{"assistant_id": assistantID})
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Supported commands are chat and init")
		os.Exit(1)
	}
}
