package registry

import (
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/containerdiscovery"
	"github.com/MustWin/cmeter/pipeline/messages/statechange"
)

const NAME = "registry"

type Filter struct {
	TrackingLabel string
	registry      *containers.Registry
}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx *pipeline.Context, m pipeline.Message) error {
	switch m.Type() {
	case containerdiscovery.TYPE:
		container := m.Body().(*containers.ContainerInfo)
		if !filter.IsTrackable(container) {
			ctx.Stop()
		} else if err := filter.registry.Register(ctx, container); err != nil {
			context.GetLogger(ctx).Errorf("error registering container: %v", err)
		}

	case statechange.TYPE:
		details := m.Body().(*statechange.Details)
		if details.State == containers.StateRunning {
			ctx.Pipeline.Send(ctx, containerdiscovery.NewMessage(details.Container))
		}

		if !filter.registry.IsRegistered(details.ContainerName) {
			ctx.Stop()
		} else if details.State == containers.StateStopped {
			if err := filter.registry.Drop(ctx, details.ContainerName); err != nil {
				context.GetLogger(ctx).Errorf("error dropping container: %v", err)
			}
		}
	}

	return nil
}

func (filter *Filter) IsTrackable(info *containers.ContainerInfo) bool {
	for k, _ := range info.Labels {
		if k == filter.TrackingLabel {
			return true
		}
	}

	return false
}

var _ pipeline.Filter = &Filter{}

func New(registry *containers.Registry, TrackingLabel string) *Filter {
	return &Filter{
		TrackingLabel: TrackingLabel,
		registry:      registry,
	}
}
