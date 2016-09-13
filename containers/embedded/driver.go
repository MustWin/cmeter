package embedded

import (
	"net/http"
	"time"

	"github.com/google/cadvisor/cache/memory"
	cadvisorMetrics "github.com/google/cadvisor/container"
	"github.com/google/cadvisor/events"
	"github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/manager"
	"github.com/google/cadvisor/utils/sysfs"

	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/containers/cadvisor"
	"github.com/MustWin/cmeter/containers/factory"
	"github.com/MustWin/cmeter/context"
)

const statsCacheDuration = 2 * time.Minute
const maxHousekeepingInterval = 15 * time.Second
const defaultHousekeepingInterval = 10 * time.Second
const allowDynamicHousekeeping = true

func init() {
	factory.Register("embedded", &driverFactory{})
}

type driverFactory struct {
}

func (factory *driverFactory) Create(parameters map[string]interface{}) (containers.Driver, error) {
}

type driver struct {
	manager manager.Manager
}

func New() containers.Driver {
	sysFs, err := sysfs.NewRealSysFs()
	if err != nil {
		return nil, err
	}

	// Create and start the cAdvisor container manager.
	m, err := manager.New(memory.New(statsCacheDuration, nil), sysFs, maxHousekeepingInterval, allowDynamicHousekeeping, cadvisorMetrics.MetricSet{cadvisorMetrics.NetworkTcpUsageMetrics: struct{}{}}, http.DefaultClient)
	if err != nil {
		return nil, err
	}

	d := &driver{
		manager: m,
	}

	if err = m.Start(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (d *driver) WatchEvents(types ...EventType) (EventsChannel, error) {
	r := events.NewRequest()
	for _, t := range types {
		r.EventType[v1.EventType(string(t))] = true
	}

	return d.manager.WatchForEvents(r)
}

func (d *driver) GetContainers() ([]*containers.ContainerInfo, error) {
	q := &v1.ContainerInfoRequest{}
	containers, err := d.manager.AllDockerContainers(q)
	if err != nil {
		return nil, err
	}

	result := make([]*containers.ContainerInfo, 0)
	for name, info := range containers {
		info := &v1.ContainerInfo{}
		info.
		localInfo := &containers.ContainerInfo{
			Name:   info.Name,
			Labels: info.Labels,
		}

		result = append(result, localInfo)
	}

	return result
}
