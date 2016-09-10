package monitor

import ()

type EventType string

const (
	EVENTS_CREATION EventType = "creation"
	EVENTS_DELETION EventType = "deletion"
)

type Event struct {
	Type        EventType
	ContainerID string
	Labels      map[string]string
	Timestamp   int64
}

type EventChannel chan *Event

type Monitor interface {
	OpenEventChannel(types ...EventType) (EventChannel, error)
	GetContainers() []string
}
