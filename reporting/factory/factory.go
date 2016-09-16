package factory

import (
	"fmt"

	"github.com/MustWin/cmeter/reporting"
)

var reportingFactories = make(map[string]ReportingDriverFactory)

type ReportingDriverFactory interface {
	Create(parameters map[string]interface{}) (reporting.Driver, error)
}

func Register(name string, factory ReportingDriverFactory) {
	if factory == nil {
		panic("ReportingDriverFactory cannot be nil")
	}

	if _, registered := reportingFactories[name]; registered {
		panic(fmt.Sprintf("ReportingDriverFactory named %s already registered", name))
	}

	reportingFactories[name] = factory
}

func Create(name string, parameters map[string]interface{}) (containers.Driver, error) {
	if factory, ok := reportingFactories[name]; ok {
		return factory.Create(parameters)
	}

	return nil, InvalidReportingDriverError{name}
}

type InvalidReportingDriverError struct {
	Name string
}

func (err InvalidReportingDriverError) Error() string {
	return fmt.Sprintf("Reporting driver not registered: %s", err.Name)
}
