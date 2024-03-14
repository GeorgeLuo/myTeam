package main

import (
	"bufio"
	"flag"
	"fmt"
	"myTeam/pkg/courier"
	"myTeam/pkg/llmclient/openai"
	"myTeam/pkg/messagebuilder"
	"myTeam/pkg/workspace"
	"os"
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

			_, messages := courier.GetMessagesByRecipient(*senderID)
			for _, message := range messages {
				fmt.Printf("%+v\n", message.Message)
			}

			// Check if the user wants to quit
			if input == "/quit" {
				fmt.Println("Exiting chat...")
				break
			}
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Supported command is chat")
		os.Exit(1)
	}
}
