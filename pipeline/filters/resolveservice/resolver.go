package resolveservice

import (
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/statechange"
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
	switch m.Type() {
	case statechange.TYPE:
		details := m.Body().(*statechange.Details)
		context.GetLoggerWithField(ctx, "container.name", details.ContainerName).Infof("state => %s", details.State)
	}

	return nil
}

func New(registry *containers.Registry, serviceKeyLabel string) pipeline.Filter {
	return &filter{
		registry:        registry,
		serviceKeyLabel: serviceKeyLabel,
	}
}
