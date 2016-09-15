package embedded

import (
	"github.com/google/cadvisor/manager"
	"sync"

	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/context"
)

type statsChannel struct {
	startFetch sync.Once
	manager    manager.Manager
	container  *containers.ContainerInfo
	ch         chan *containers.Stats
	closer     chan bool
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
				close(ch.ch)
				return
			} else {
				ch.manager.// TODO: make call here
			}
		}
	}
}

func (ch *statsChannel) Close() error {
	ch.closer <- true
	return nil
}

func newStatsChannel(manager manager.Manager, container *containers.ContainerInfo) *statsChannel {
	return &statsChannel{
		manager:   manager,
		container: container,
	}
}
