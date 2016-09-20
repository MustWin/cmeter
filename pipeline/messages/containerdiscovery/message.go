package containerdiscovery

import (
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/pipeline"
)

const TYPE = "container_discovery"

type Message struct {
	Container *containers.ContainerInfo
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.Container
}

var _ pipeline.Message = &Message{}

func NewMessage(container *containers.ContainerInfo) *Message {
	return &Message{
		Container: container,
	}
}
