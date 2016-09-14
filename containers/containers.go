package containers

import (
	"time"

	"github.com/MustWin/cmeter/context"
)

type EventType string

const (
	EventContainerCreation EventType = "containerCreation"
	EventContainerDeletion EventType = "containerDeletion"
	EventContainerOom      EventType = "oom"
	EventContainerOomKill  EventType = "oomKill"
)

type State string

const (
	StateRunning State = "running"
	StateStopped State = "stopped"
	StateUnknown State = "unknown"
)

type Event struct {
	Type          EventType
	ContainerName string
	ServiceKey    string
	Timestamp     time.Time
}

type EventsChannel interface {
	GetChannel() <-chan *Event
	Close() error
}

type ContainerInfo struct {
	Name   string
	Labels map[string]string
}

type Metrics struct {
}

type MetricsChannel interface {
	Container() *ContainerInfo
	GetChannel() <-chan *Metrics
	Close() error
}

type Driver interface {
	WatchEvents(ctx context.Context, types ...EventType) (EventsChannel, error)
	GetContainers(ctx context.Context) ([]*ContainerInfo, error)
	GetContainer(ctx context.Context, name string) (*ContainerInfo, error)
	GetContainerMetrics(ctx context.Context, container *containers.ContainerInfo) (MetricsChannel, error)
	CloseAllChannels(ctx context.Context) error
}
