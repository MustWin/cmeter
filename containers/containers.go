package containers

import (
	"github.com/MustWin/cmeter/context"
)

type EventType string

const (
	EventContainerCreation EventType = "containerCreation"
	EventContainerDeletion EventType = "containerDeletion"
	EventContainerOom      EventType = "oom"
	EventContainerOomKill  EventType = "oomKill"
	EventContainerExisted  EventType = "containerExisted"
)

type State string

const (
	StateRunning State = "running"
	StateStopped State = "stopped"
	StateUnknown State = "unknown"
)

func StateFromEvent(eventType EventType) State {
	switch eventType {
	case EventContainerCreation:
		return StateRunning
	case EventContainerDeletion:
		return StateStopped
	}

	return StateUnknown
}

type StateChange struct {
	State     State          `json:"state"`
	Source    *Event         `json:"source_event"`
	Container *ContainerInfo `json:"container"`
}

type Event struct {
	Type          EventType `json:"type"`
	ContainerName string    `json:""`
	Timestamp     int64     `json:"timestamp"`
}

type EventsChannel interface {
	GetChannel() <-chan *Event
	Close() error
}

type ContainerInfo struct {
	Name   string
	Labels map[string]string
}

type Stats struct {
}

type StatsChannel interface {
	Container() *ContainerInfo
	GetChannel() <-chan *Stats
	Close() error
}

type Driver interface {
	WatchEvents(ctx context.Context, types ...EventType) (EventsChannel, error)
	GetContainers(ctx context.Context) ([]*ContainerInfo, error)
	GetContainer(ctx context.Context, name string) (*ContainerInfo, error)
	GetContainerStats(ctx context.Context, container *ContainerInfo) (StatsChannel, error)
	CloseAllChannels(ctx context.Context) error
}
