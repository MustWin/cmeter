package uplink

import (
	"fmt"

	"github.com/MustWin/cmeter/api"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/transmit"
)

const NAME = "uplink"

type Filter struct {
	client api.Client
}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx *pipeline.Context, m pipeline.Message) error {
	switch m.Type() {
	case transmit.TYPE:
		event := m.Body().(*api.Event)
		if err := filter.client.Send(event); err != nil {
			return fmt.Errorf("couldn't send event: %v", err)
		}

		ctx.Stop()
	}

	return nil
}

func New(client api.Client) *Filter {
	return &Filter{
		client: client,
	}
}
