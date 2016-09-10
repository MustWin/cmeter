package statechange

import (
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/monitor"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/shared/uuid"
)

const TYPE = "state_change"

type Message struct {
	id          string
	ContainerID string
	State       string
}

func (msg *Message) ID() string {
	return msg.id
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.ContainerID
}

func NewMessage(containerId string) pipeline.Message {
	return &Message{
		id:          uuid.Generate(),
		ContainerID: containerId,
	}
}
