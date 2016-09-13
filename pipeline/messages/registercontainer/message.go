package registercontainer

import (
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/shared/uuid"
)

const TYPE = "register_container"

type Message struct {
	id        string
	Container *containers.ContainerInfo
}

func (msg *Message) ID() string {
	return msg.id
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.Container
}

func NewMessage(container *containers.ContainerInfo) pipeline.Message {
	return &Message{
		id:        uuid.Generate(),
		Container: container,
	}
}
