package embedded

import (
	"github.com/google/cadvisor/events"

	"github.com/MustWin/cmeter/containers"
)

type eventChannel struct {
	inner   *events.EventChannel
	channel chan *containers.Event
}

func newEventChannel(cec *events.EventChannel) *eventChannel {
	ec := &eventChannel{
		inner:   cec,
		channel: make(chan *containers.Event),
	}

	go func() {
		defer close(ec.channel)
		for src := range cec.GetChannel() {
			e := &containers.Event{
				Container: &containers.ContainerInfo{
					Name: src.ContainerName,
				},
				Timestamp: src.Timestamp.Unix(),
				Type:      containers.EventType(string(src.EventType)),
			}

			ec.channel <- e
		}
	}()

	return ec
}

func (ec *eventChannel) GetChannel() <-chan *containers.Event {
	return ec.channel
}

func (ec *eventChannel) Close() error {
	// TODO: may not need
	return nil
}
