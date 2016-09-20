package statechange

import (
	"time"

	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/base"
)

const TYPE = "state_change"

type Details struct {
	ContainerName string
	Container     *containers.ContainerInfo
	State         containers.State
	Timestamp     time.Time
}

type Message struct {
	*base.Message
	Details *Details
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.Details
}

var _ pipeline.Message = &Message{}

func NewMessage(event *containers.Event) *Message {
	details := &Details{
		ContainerName: event.ContainerName,
		Timestamp:     event.Timestamp,
	}

	switch event.Type {
	case containers.EventContainerCreation:
		details.State = containers.StateRunning
	case containers.EventContainerDeletion:
		details.State = containers.StateStopped
	default:
		details.State = containers.StateUnknown
	}

	return &Message{
		Message: &base.Message{},
		Details: details,
	}
}
