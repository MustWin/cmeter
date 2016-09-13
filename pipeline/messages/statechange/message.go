package statechange

import (
	"time"

	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/shared/uuid"
)

const TYPE = "state_change"

type Details struct {
	ContainerName string
	State         string
	Timestamp     time.Time
}

type Message struct {
	id      string
	Details *Details
}

func (msg *Message) ID() string {
	return msg.id
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.Details
}

func NewMessage(event *containers.Event) pipeline.Message {
	details := &Details{
		ContainerName: event.ContainerName,
		Timestamp:     event.Timestamp,
	}

	switch event.Type {
	case containers.EventContainerCreation:
		details.State = "created"
	case containers.EventContainerDeletion:
		details.State = "deleted"
	}

	return &Message{
		id:      uuid.Generate(),
		Details: details,
	}
}
