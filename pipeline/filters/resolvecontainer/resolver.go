package resolvecontainer

import (
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/statechange"
)

const NAME = "container_resolver"

type filter struct {
	containers containers.Driver
}

func (filter *filter) Name() string {
	return NAME
}

func (filter *filter) HandleMessage(ctx *pipeline.Context, m pipeline.Message) error {
	switch m.Type() {
	case statechange.TYPE:
		details := m.Body().(*statechange.Details)
		if details.Container == nil {
			info, err := filter.containers.GetContainer(ctx, details.ContainerName)
			if err != nil {
				return err
			}

			details.Container = info
		}
	}

	return nil
}

func New(driver containers.Driver) pipeline.Filter {
	return &filter{
		containers: driver,
	}
}
