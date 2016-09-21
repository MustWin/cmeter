package resolvecontainer

import (
	"time"

	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/statechange"
)

const NAME = "container_resolver"

type Filter struct {
	containers containers.Driver
	registry   *containers.Registry
}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx context.Context, m pipeline.Message) error {
	switch m.Type() {
	case statechange.TYPE:
		details := m.Body().(*statechange.Details)
		if details.Container == nil {
			switch details.State {
			case containers.StateRunning:
				time.Sleep(time.Millisecond * 5000)
				info, err := filter.containers.GetContainer(ctx, details.ContainerName)
				if err != nil {
					if err == containers.ErrContainerNotFound {
						return nil
					}

					return err
				}

				details.Container = info
			case containers.StateStopped:
				if filter.registry.IsRegistered(details.ContainerName) {
					details.Container, _ = filter.registry.Get(details.ContainerName)
				}
			}
		}
	}

	return nil
}

var _ pipeline.Filter = &Filter{}

func New(driver containers.Driver, registry *containers.Registry) *Filter {
	return &Filter{
		containers: driver,
		registry:   registry,
	}
}
