package containersample

import (
	"github.com/MustWin/cmeter/collector"
	"github.com/MustWin/cmeter/pipeline"
)

const TYPE = "container_sample"

type Message struct {
	Sample *collector.Sample
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.Sample
}

var _ pipeline.Message = &Message{}

func NewMessage(sample *collector.Sample) *Message {
	return &Message{
		Sample: sample,
	}
}
