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
	ch     containers.UsageChannel
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
	Usage     *containers.Usage         `json:"usage"`
}

type MachineSample struct {
	Timestamp int64                    `json:"timestamp"`
	FrameSize time.Duration            `json:"rate"`
	Machine   *containers.MachineInfo  `json:"machine"`
	Usage     *containers.MachineUsage `json:"usage"`
}

func (c *Collector) Num() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.collections)
}

func (c *Collector) Collect(ctx context.Context, ch containers.UsageChannel) error {
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
		case usage, ok := <-data.ch.GetChannel():
			if !ok {
				defer context.GetLogger(ctx).Info("container stats collection completed")
				if _, err := c.Stop(ctx, data.ch.Container()); err != nil {
					context.GetLogger(ctx).Errorf("error stopping container stats collection: %v", err)
				}

				return
			} else {
				sample := &Sample{
					Container: data.ch.Container(),
					Usage:     usage,
					Timestamp: time.Now().Unix(),
					FrameSize: c.Rate,
				}

				c.samples <- sample
			}
		}
	}
}

func (c *Collector) Stop(ctx context.Context, container *containers.ContainerInfo) (containers.UsageChannel, error) {
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

func (c *Collector) StopAll() ([]containers.UsageChannel, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	channels := make([]containers.UsageChannel, 0)
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

type MachineCollector struct {
	context.Context
	feed    containers.MachineUsageFeed
	Rate    time.Duration
	active  bool
	mutex   sync.Mutex
	samples chan *MachineSample
}

func NewMachine(ctx context.Context, feed containers.MachineUsageFeed, config configuration.CollectorConfig) *MachineCollector {
	return &MachineCollector{
		Context: ctx,
		Rate:    time.Duration(config.Rate) * time.Millisecond,
		feed:    feed,
		active:  false,
		samples: make(chan *MachineSample, CHANNEL_BUFFER_SIZE),
	}
}

func (c *MachineCollector) Active() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.active
}

func (c *MachineCollector) Start() error {
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

func (c *MachineCollector) GetChannel() <-chan *MachineSample {
	return c.samples
}

func (c *MachineCollector) Stop() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.active = false
	context.GetLogger(c).Info("stopped machine stats collection")
	return nil
}

func (c *MachineCollector) doCollect() {
	t := time.NewTicker(c.Rate)
	for _ = range t.C {
		if !c.Active() {
			t.Stop()
			return
		}

		usage := c.feed.Next()
		if usage == nil {
			context.GetLogger(c).Error("couldn't sample machine stats")
			continue
		}

		sample := &MachineSample{
			Machine:   c.feed.Machine(),
			Usage:     usage,
			FrameSize: c.Rate,
			Timestamp: time.Now().Unix(),
		}

		c.samples <- sample
	}
}
