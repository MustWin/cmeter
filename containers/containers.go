package containers

import (
	"time"
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
)

type Event struct {
	Type          EventType
	ContainerName string
	ServiceKey    string
	Timestamp     time.Time
}

type EventsChannel interface {
	GetChannel() <-chan *Event
}

type ContainerInfo struct {
	Name   string
	Labels map[string]string
}

type Driver interface {
	WatchEvents(types ...EventType) (EventsChannel, error)
	GetContainers() ([]*ContainerInfo, error)
}
