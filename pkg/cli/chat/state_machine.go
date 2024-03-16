package chat

import "fmt"

type State int

const (
	AwaitingInput State = iota
	SendingMessage
	ReceivingResponse
	ReviewingAttachments
	Exiting
)

type Event int

const (
	MessageEntered Event = iota
	SendMessage
	MessageSent
	ReceivedResponse
	ReviewAttachmentsInitiated
	AttachmentsReviewed
	QuitCommandEntered
)

type ChatStateMachine struct {
	CurrentState State
	Transitions  map[State]map[Event]State
}

func NewChatStateMachine() *ChatStateMachine {
	transitions := map[State]map[Event]State{
		AwaitingInput: {
			MessageEntered:     SendingMessage,
			QuitCommandEntered: Exiting,
		},
		SendingMessage: {
			MessageSent: ReceivingResponse,
		},
		ReceivingResponse: {
			ReceivedResponse:           AwaitingInput,
			ReviewAttachmentsInitiated: ReviewingAttachments,
		},
		ReviewingAttachments: {
			AttachmentsReviewed: AwaitingInput,
		},
	}
	return &ChatStateMachine{
		CurrentState: AwaitingInput,
		Transitions:  transitions,
	}
}

func (sm *ChatStateMachine) Transition(event Event) error {
	if nextState, ok := sm.Transitions[sm.CurrentState][event]; ok {
		sm.CurrentState = nextState
		return nil
	}
	return fmt.Errorf("invalid transition")
}

func (sm *ChatStateMachine) CanTransition(event Event) bool {
	_, ok := sm.Transitions[sm.CurrentState][event]
	return ok
}
