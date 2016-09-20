package devnull

import (
	"github.com/MustWin/cmeter/pipeline"
)

const TYPE = "devnull"

type Message struct{}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return nil
}

var _ pipeline.Message = &Message{}
