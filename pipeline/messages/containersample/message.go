package containersample

import (
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/shared/uuid"
)

const TYPE = "container_sample"

type Message struct {
	id        string
	Container *containers.ContainerInfo
	Sample    *containers.Metrics
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

func NewMessage(container *containers.ContainerInfo, sample *containers.Metrics) pipeline.Message {
	return &Message{
		id:        uuid.Generate(),
		Container: container,
		Sample:    sample,
	}
}
