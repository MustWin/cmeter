package reportevent

import (
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/reporting"
)

const TYPE = "report_event"

type Message struct {
	Event *reporting.Event
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.Event
}

var _ pipeline.Message = &Message{}

func NewMessage(event *reporting.Event) *Message {
	return &Message{
		Event: event,
	}
}
