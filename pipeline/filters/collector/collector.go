package collector

import (
	"github.com/MustWin/cmeter/collector"
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/containerdiscovery"
	"github.com/MustWin/cmeter/pipeline/messages/statechange"
)

const NAME = "sample_collector"

type Filter struct {
	containers containers.Driver
	collector  *collector.Collector
}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx context.Context, m pipeline.Message) error {
	var container *containers.ContainerReference
	drop := false

	switch m.Type() {
	case containerdiscovery.TYPE:
		info := m.Body().(*containers.ContainerInfo)
		container = info.ContainerReference

	case statechange.TYPE:
		change := m.Body().(*containers.StateChange)
		container = change.Container.ContainerReference
		if change.State == containers.StateStopped {
			drop = true
		}
	}

	var err error
	if container != nil {
		if drop {
			_, err = filter.collector.Stop(ctx, container)
		} else {
			var ch containers.StatsChannel
			ch, err = filter.containers.GetContainerStats(ctx, container)
			if err == nil {
				err = filter.collector.Collect(ctx, ch)
			}
		}
	}

	return err
}

var _ pipeline.Filter = &Filter{}

func New(driver containers.Driver, collector *collector.Collector) *Filter {
	return &Filter{
		containers: driver,
		collector:  collector,
	}
}
