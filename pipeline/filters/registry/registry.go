package registry

import (
	"fmt"
	"time"

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

func (filter *Filter) HandleMessage(ctx context.Context, m pipeline.Message) error {
	switch m.Type() {
	case containerdiscovery.TYPE:
		container := m.Body().(*containers.ContainerInfo)
		if filter.IsTrackable(container) {
			if err := filter.registry.Register(ctx, container); err != nil {
				return fmt.Errorf("error registering container: %v", err)
			}

			return pipeline.SetMessage(ctx, statechange.NewMessage(&containers.StateChange{
				State: containers.StateRunning,
				Source: &containers.Event{
					Type:      containers.EventContainerExisted,
					Timestamp: time.Now().Unix(),
				},
				Container: container,
			}))
		} else {
			return pipeline.StopProcessing(ctx)
		}

	case statechange.TYPE:
		change := m.Body().(*containers.StateChange)
		if !filter.registry.IsRegistered(change.Container.Name) {
			return pipeline.StopProcessing(ctx)
		} else if change.State == containers.StateStopped {
			if err := filter.registry.Drop(ctx, change.Container.Name); err != nil {
				return fmt.Errorf("error dropping container: %v", err)
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
