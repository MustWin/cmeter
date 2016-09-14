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
	trackingLabel string
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
			context.GetLogger(ctx)
		}

	case statechange.TYPE:
		details := m.Body().(*statechange.Details)
		if details.State == containers.StateRunning {
			if err := ctx.Pipeline.Send(ctx, containerdiscovery.NewMessage(details.Container)); err != nil {
				return err
			}
		}

		if !filter.registry.IsRegistered(details.ContainerName) {
			ctx.Stop()
		}
	}

	return nil
}

func (filter *Filter) IsTrackable(info *containers.ContainerInfo) bool {
	for k, _ := range info.Labels {
		if k == filter.trackingLabel {
			return true
		}
	}

	return false
}

func New(registry *containers.Registry, trackingLabel string) *Filter {
	return &Filter{
		trackingLabel: trackingLabel,
		registry:      registry,
	}
}
