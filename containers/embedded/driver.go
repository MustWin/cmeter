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
	"github.com/MustWin/cmeter/containers/factory"
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

	return d, nil
}

type driver struct {
	manager manager.Manager
}

func (d *driver) WatchEvents(types ...containers.EventType) (containers.EventsChannel, error) {
	r := events.NewRequest()
	for _, t := range types {
		r.EventType[v1.EventType(string(t))] = true
	}

	cec, err := d.manager.WatchForEvents(r)
	if err != nil {
		return nil, err
	}

	return newEventChannel(cec), nil
}

func (d *driver) GetContainers() ([]*containers.ContainerInfo, error) {
	q := &v1.ContainerInfoRequest{}
	rawContainers, err := d.manager.AllDockerContainers(q)
	if err != nil {
		return nil, err
	}

	result := make([]*containers.ContainerInfo, 0)
	for _, info := range rawContainers {
		localInfo := &containers.ContainerInfo{
			Name:   info.Name,
			Labels: info.Labels,
		}

		result = append(result, localInfo)
	}

	return result, nil
}
