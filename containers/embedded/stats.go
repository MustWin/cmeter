package embedded

import (
	"errors"
	"sync"

	"github.com/google/cadvisor/manager"

	"github.com/MustWin/cmeter/containers"
)

type statsChannel struct {
	startFetch  sync.Once
	manager     manager.Manager
	container   *containers.ContainerInfo
	ch          chan *containers.Stats
	closer      chan bool
	closedMutex sync.Mutex
	closed      bool
}

func (ch *statsChannel) Container() *containers.ContainerInfo {
	return ch.container
}

func (ch *statsChannel) GetChannel() <-chan *containers.Stats {
	ch.startFetch.Do(func() {
		go ch.startChannel()
	})

	return ch.ch
}

func (ch *statsChannel) startChannel() {
	for {
		select {
		case done, ok := <-ch.closer:
			if !ok && done {
				ch.closedMutex.Lock()
				defer ch.closedMutex.Unlock()
				close(ch.ch)
				ch.closed = true
				return
			} else {
				//ch.manager.
			}
		}
	}
}

func (ch *statsChannel) Close() error {
	ch.closedMutex.Lock()
	defer ch.closedMutex.Unlock()
	if ch.closed {
		return errors.New("already closed")
	}

	ch.closer <- true
	return nil
}

func newStatsChannel(manager manager.Manager, container *containers.ContainerInfo) *statsChannel {
	return &statsChannel{
		manager:   manager,
		container: container,
		closed:    false,
	}
}
