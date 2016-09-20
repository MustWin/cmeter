package base

import (
	"sync"

	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/shared/uuid"
)

const TYPE = "base"

type Message struct {
	id    string
	idgen sync.Once
}

func (msg *Message) ID() string {
	msg.idgen.Do(func() {
		msg.id = uuid.Generate()
	})

	return msg.id
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return nil
}

var _ pipeline.Message = &Message{}
