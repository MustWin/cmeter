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
	ch     containers.StatsChannel
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
	Timestamp int64                     `json:"timestamp"`
	FrameSize time.Duration             `json:"rate"`
	Container *containers.ContainerInfo `json:"container"`
	Stats     *containers.Stats         `json:"stats"`
}

type HostSample struct {
	Timestamp int64                    `json:"timestamp"`
	FrameSize time.Duration            `json:"rate"`
	Machine   *containers.MachineInfo  `json:"machine"`
	Stats     *containers.MachineStats `json:"stats"`
}

func (c *Collector) Num() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.collections)
}

func (c *Collector) Collect(ctx context.Context, ch containers.StatsChannel) error {
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
	log := context.GetLoggerWithField(ctx, "container.name", ch.Container().Name)
	cctx := context.WithLogger(ctx, log)
	go c.doCollect(cctx, data)
	log.Info("started container stats collection")
	return nil
}

func (c *Collector) GetChannel() <-chan *Sample {
	return c.samples
}

func (c *Collector) doCollect(ctx context.Context, data *collectorData) {
	for _ = range data.ticker.C {
		select {
		case metrics, ok := <-data.ch.GetChannel():
			if !ok {
				defer context.GetLogger(ctx).Info("container stats collection completed")
				if _, err := c.Stop(ctx, data.ch.Container()); err != nil {
					context.GetLogger(ctx).Errorf("error stopping container stats collection: %v", err)
				}

				return
			} else {
				sample := &Sample{
					Container: data.ch.Container(),
					Stats:     metrics,
					Timestamp: time.Now().Unix(),
					FrameSize: c.Rate,
				}

				c.samples <- sample
			}
		}
	}
}

func (c *Collector) Stop(ctx context.Context, container *containers.ContainerInfo) (containers.StatsChannel, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	data, ok := c.collections[container.Name]
	if !ok {
		return nil, fmt.Errorf("no collection for %s", container.Name)
	}

	data.ticker.Stop()
	delete(c.collections, container.Name)
	context.GetLoggerWithField(ctx, "container.name", data.ch.Container().Name).Info("stopped container stats collection")
	return data.ch, nil
}

func (c *Collector) StopAll() ([]containers.StatsChannel, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	channels := make([]containers.StatsChannel, 0)
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

type HostCollector struct {
	context.Context
	feed    containers.MachineStatsFeed
	Rate    time.Duration
	active  bool
	mutex   sync.Mutex
	samples chan *HostSample
}

func NewHost(ctx context.Context, feed containers.MachineStatsFeed, config configuration.CollectorConfig) *HostCollector {
	return &HostCollector{
		Context: ctx,
		Rate:    time.Duration(config.Rate) * time.Millisecond,
		feed:    feed,
		active:  false,
		samples: make(chan *HostSample, CHANNEL_BUFFER_SIZE),
	}
}

func (c *HostCollector) Active() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.active
}

func (c *HostCollector) Start() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.active {
		return fmt.Errorf("machine stats collector already started")
	}

	c.active = true
	go c.doCollect()
	context.GetLogger(c).Info("started machine stats collection")
	return nil
}

func (c *HostCollector) GetChannel() <-chan *HostSample {
	return c.samples
}

func (c *HostCollector) Stop() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.active = false
	context.GetLogger(c).Info("stopped machine stats collection")
	return nil
}

func (c *HostCollector) doCollect() {
	t := time.NewTicker(c.Rate)
	for _ = range t.C {
		if !c.Active() {
			t.Stop()
			return
		}

		metrics := c.feed.Next()
		if metrics == nil {
			context.GetLogger(c).Error("couldn't sample machine stats")
			continue
		}

		sample := &HostSample{
			Machine:   c.feed.Machine(),
			Stats:     metrics,
			FrameSize: c.Rate,
			Timestamp: time.Now().Unix(),
		}

		c.samples <- sample
	}
}
