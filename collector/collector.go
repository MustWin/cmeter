package collector

import (
	"fmt"
	"time"

	"github.com/MustWin/cmeter/"
	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/containersample"
)

type collectorData struct {
	ch     containers.MetricChannel
	ticker *time.Ticker
}

type Collector struct {
	context.Context
	Pipeline    pipeline.Pipeline
	Rate        time.Duration
	collections map[string]*collectorData
}

func (c *Collector) Collect(ch containers.MetricsChannel) error {
	if c.collections == nil {
		c.collections = make(map[string]*collectorData)
	}

	data := &collectorData{
		ch:     ch,
		ticker: time.NewTicker(c.Rate),
	}

	c.collections[ch.Channel().Name] = data
	// log
	go c.doCollect(data)
	return nil
}

func (c *Collector) doCollect(data *collectorData) {
	for _ = range data.ticker.C {
		select {
		case sample, ok := <-data.ch.GetChannel():
			if !ok {
				if err := c.Stop(data.ch.Container()); err != nil {
					//log
				}
			} else {
				m := containersample.NewMessage(data.ch.Container(), sample)
				c.Pipeline.Send(c, m)
			}
		}
	}
}

func (c *Collector) Stop(container *containers.ContainerInfo) (containers.MetricsChannel, error) {
	data, ok := c.collections[container.Name]
	if !ok {
		return nil, fmt.Errorf("no collection for %s", container.Name)
	}

	data.ticker.Stop()
	delete(c.collections, container.Name)
	return
}

func (c *Collector) StopAll() ([]containers.MetricsChannel, error) {
	channels := make([]containers.MetricsChannel, 0)
	for _, data := range c.collections {
		data.ticker.Stop()
		channels = append(channels, data.ch)
	}

	c.collections = make(map[string]*collectorData)
	return channels, nil
}

func New(ctx context.Context, config *configuration.CollectorConfig, pipeline pipeline.Pipeline) *Collector {
	return &Collector{
		Context:  ctx,
		Pipeline: pipeline,
		Rate:     time.Duration(config.Rate * time.Millisecond),
	}
}
