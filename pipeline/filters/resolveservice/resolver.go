package resolveservice

import (
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
)

const NAME = "resolver"

type filter struct {
	serviceKeyLabel string
	registry        *containers.Registry
}

func (filter *filter) Name() string {
	return NAME
}

func (filter *filter) HandleMessage(ctx *pipeline.Context, m pipeline.Message) error {
	context.GetLogger(ctx).Infof("processing %q message", m.Type())
	return nil
}

func New(registry *containers.Registry, serviceKeyLabel string) pipeline.Filter {
	return &filter{
		registry:        registry,
		serviceKeyLabel: serviceKeyLabel,
	}
}
