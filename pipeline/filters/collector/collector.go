package collector

import (
	"github.com/MustWin/cmeter/collector"
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/containerdiscovery"
	"github.com/MustWin/cmeter/pipeline/messages/statechange"
)

const NAME = "sample_collector"

type filter struct {
	containers containers.Driver
	collector  *collector.Collector
}

func (filter *filter) Name() string {
	return NAME
}

func (filter *filter) HandleMessage(ctx *pipeline.Context, m pipeline.Message) error {
	var container *containers.ContainerInfo
	drop := false

	switch m.Type() {
	case containerdiscovery.TYPE:
		container = m.Body().(*containers.ContainerInfo)

	case statechange.TYPE:
		details := m.Body().(*statechange.Details)
		container = details.Container
		if details.State != containers.StateRunning {
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

func New(driver containers.Driver, collector *collector.Collector) pipeline.Filter {
	return &filter{
		containers: driver,
		collector:  collector,
	}
}
