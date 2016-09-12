package containers

import ()

type EventType string

const (
	EventContainerCreation = "creation"
	EventContainerDeletion = "deletion"
)

type Event struct {
	Type          EventType
	ContainerName string
	Labels        map[string]string
	Timestamp     time.time
}

type EventsChannel interface {
	GetChannels() chan *Event
}

type Driver interface {
	WatchEvents(types ...EventType) (EventsChannel, error)
	GetContainers() ([]string, error)
}
