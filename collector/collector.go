package collector

import (
	"fmt"
	"sync"
	"time"

	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/context"
)

const CHANNEL_BUFFER_SIZE = 3000

type collectorData struct {
	ch     containers.MetricsChannel
	ticker *time.Ticker
}

type Collector struct {
	context.Context
	Rate        time.Duration
	collections map[string]*collectorData
	samples     chan *Sample
	mutex       sync.Mutex
}

type Sample struct {
	Container *containers.ContainerInfo
	Metrics   *containers.Metrics
}

func (c *Collector) Collect(ch containers.MetricsChannel) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.collections == nil {
		c.collections = make(map[string]*collectorData)
	}

	data := &collectorData{
		ch:     ch,
		ticker: time.NewTicker(c.Rate),
	}

	c.collections[ch.Container().Name] = data
	// log
	go c.doCollect(data)
	return nil
}

func (c *Collector) GetChannel() <-chan *Sample {
	return c.samples
}

func (c *Collector) doCollect(data *collectorData) {
	for _ = range data.ticker.C {
		select {
		case metrics, ok := <-data.ch.GetChannel():
			if !ok {
				if _, err := c.Stop(data.ch.Container()); err != nil {
					//log
				}
			} else {
				sample := &Sample{
					Container: data.ch.Container(),
					Metrics:   metrics,
				}

				c.samples <- sample
			}
		}
	}
}

func (c *Collector) Stop(container *containers.ContainerInfo) (containers.MetricsChannel, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	data, ok := c.collections[container.Name]
	if !ok {
		return nil, fmt.Errorf("no collection for %s", container.Name)
	}

	data.ticker.Stop()
	delete(c.collections, container.Name)
	return data.ch, nil
}

func (c *Collector) StopAll() ([]containers.MetricsChannel, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	channels := make([]containers.MetricsChannel, 0)
	for _, data := range c.collections {
		data.ticker.Stop()
		channels = append(channels, data.ch)
	}

	c.collections = make(map[string]*collectorData)
	return channels, nil
}

func New(config configuration.CollectorConfig) *Collector {
	return &Collector{
		Rate:    time.Duration(config.Rate) * time.Millisecond,
		samples: make(chan *Sample, CHANNEL_BUFFER_SIZE),
	}
}
