package containerdiscovery

import (
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/base"
)

const TYPE = "container_discovery"

type Message struct {
	*base.Message
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
		Message:   &base.Message{},
		Container: container,
	}
}
