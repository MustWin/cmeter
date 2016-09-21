package statechange

import (
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/pipeline"
)

const TYPE = "state_change"

type Message struct {
	body *containers.StateChange
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.body
}

var _ pipeline.Message = &Message{}

func NewMessage(change *containers.StateChange) *Message {
	return &Message{
		body: bod
	}
	body := &containers.StateChange{
		Source:    event,
		Container: nil,
	}

	switch event.Type {
	case containers.EventContainerCreation:
		body.State = containers.StateRunning
	case containers.EventContainerDeletion:
		body.State = containers.StateStopped
	default:
		body.State = containers.StateUnknown
	}

	return &Message{
		body: body,
	}
}
