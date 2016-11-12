package disposer

import ()

// Not thread safe
type Disposer struct {
	doneCh  chan struct{}
	quitChs []chan struct{}
}

func (d *Disposer) Quitter() chan struct{} {
	ch := make(chan struct{}, 0)
	if d.quitChs == nil {
		d.quitChs = []chan struct{}{ch}
	} else {
		d.quitChs = append(d.quitChs, ch)
	}

	return ch
}

func (d *Disposer) Wait() {
	<-d.doneCh
}

func (d *Disposer) Dispose() {
	d.doneCh <- struct{}{}
}

func (d *Disposer) QuitAll() {
	v := struct{}{}
	for _, ch := range d.quitChs {
		ch <- v
	}
}

func New() *Disposer {
	return &Disposer{
		doneCh:  make(chan struct{}, 0),
		quitChs: make([]chan struct{}, 0),
	}
}
