package transmit

import (
	"github.com/MustWin/cmeter/api"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/shared/uuid"
)

const TYPE = "transmit"

type Message struct {
	id   string
	body *api.Event
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

func NewMessage(event *api.Event) pipeline.Message {
	return &Message{
		id:   uuid.Generate(),
		body: event,
	}
}
