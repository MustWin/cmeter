package reportevent

import (
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/reporting"
	"github.com/MustWin/cmeter/shared/uuid"
)

const TYPE = "report_event"

type Message struct {
	id   string
	body *reporting.Event
}

func (msg *Message) ID() string {
	return msg.id
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.body
}

func NewMessage(event *reporting.Event) pipeline.Message {
	return &Message{
		id:   uuid.Generate(),
		body: event,
	}
}
