package embedded

import (
	"github.com/google/cadvisor/manager"

	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/containers/cadvisor"
	"github.com/MustWin/cmeter/containers/factory"
	"github.com/MustWin/cmeter/context"
)

const NAME = "cadvisor"

func init() {
	factory.Register(NAME, &driverFactory{})
}

type driverFactory struct {
}

func (factory *driverFactory) Create(parameters map[string]interface{}) (containers.Driver, error) {
	return nil, nil
}

type driver struct {
	manager.Manager
	client containers.Driver
}

func New() containers.Driver {

}

func (d *driver) WatchEvents(types ...EventType) (EventsChannel, error) {
}

func (d *driver) GetContainers() ([]string, error) {

}
