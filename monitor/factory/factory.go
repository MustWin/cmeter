package factory

import (
	"fmt"

	"github.com/MustWin/cmeter/monitor"
)

var monitorFactories = make(map[string]MonitorFactory)

type MonitorFactory interface {
	Create(parameters map[string]interface{}) (monitor.Monitor, error)
}

func Register(name string, factory MonitorFactory) {
	if factory == nil {
		panic("MonitorFactory cannot be nil")
	}

	if _, registered := monitorFactories[name]; registered {
		panic(fmt.Sprintf("MonitorFactory named %s already registered", name))
	}

	monitorFactories[name] = factory
}

func Create(name string, parameters map[string]interface{}) (monitor.Monitor, error) {
	if factory, ok := monitorFactories[name]; ok {
		return factory.Create(parameters)
	}

	return nil, InvalidMonitorError{name}
}

type InvalidMonitorError struct {
	Name string
}

func (err InvalidMonitorError) Error() string {
	return fmt.Sprintf("Monitor not registered: %s", err.Name)
}
