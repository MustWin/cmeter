package cadvisor

import (
	"github.com/MustWin/cmeter/containers"
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

type Driver struct {
}

func New() containers.Driver {

}

func (driver *Driver) WatchEvents() (*containers.EventsChannel, error) {
	return nil, nil
}

func (driver *Driver) GetContainers() ([]string, error) {
	return nil, nil
}
