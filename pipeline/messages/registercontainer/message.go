package registercontainer

import (
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/shared/uuid"
)

const TYPE = "register_container"

type Message struct {
	id          string
	ContainerID string
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
