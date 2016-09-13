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
		for src := range cec.GetChannel() {
			e := &containers.Event{
				ContainerName: src.ContainerName,
				Timestamp:     src.Timestamp,
				Type:          containers.EventType(string(src.EventType)),
			}

			ec.channel <- e
		}

		close(ec.channel)
	}()

	return ec
}

func (ec *eventChannel) GetChannel() <-chan *containers.Event {
	return ec.channel
}
