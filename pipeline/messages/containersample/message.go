package containersample

import (
	"github.com/MustWin/cmeter/collector"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/shared/uuid"
)

const TYPE = "container_sample"

type Message struct {
	id     string
	Sample *collector.Sample
}

func (msg *Message) ID() string {
	return msg.id
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.Sample
}

func NewMessage(sample *collector.Sample) pipeline.Message {
	return &Message{
		id:     uuid.Generate(),
		Sample: sample,
	}
}
