package embedded

import (
	"errors"
	"flag"
	"net/http"
	"sync"
	"time"

	"github.com/google/cadvisor/cache/memory"
	cadvisorMetrics "github.com/google/cadvisor/container"
	"github.com/google/cadvisor/events"
	"github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/manager"
	"github.com/google/cadvisor/utils/sysfs"

	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/containers/factory"
	"github.com/MustWin/cmeter/context"
)

var parseOnce sync.Once

const (
	statsCacheDuration          = 2 * time.Minute
	maxHousekeepingInterval     = 15 * time.Second
	defaultHousekeepingInterval = 10 * time.Second
	allowDynamicHousekeeping    = true
)

func init() {
	factory.Register("embedded", &driverFactory{})
}

type driverFactory struct{}

func (factory *driverFactory) Create(parameters map[string]interface{}) (containers.Driver, error) {
	if !flag.Parsed() {
		parseOnce.Do(func() {
			flag.Parse()
		})
	}

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

func init() {
	// Override cAdvisor flag defaults
	flagOverrides := map[string]string{
		// Override the default cAdvisor housekeeping interval.
		"housekeeping_interval": defaultHousekeepingInterval.String(),
		// Disable event storage by default.
		"event_storage_event_limit": "default=0",
		"event_storage_age_limit":   "default=0",
	}

	for name, defaultValue := range flagOverrides {
		if f := flag.Lookup(name); f != nil {
			f.DefValue = defaultValue
			f.Value.Set(defaultValue)
			// TODO: can't log error here but marking as a *maybe* and might find a more elegant approach for this
			//log.Errorf("Expected cAdvisor flag %q not found", name)
		}
	}
}

type driver struct {
	manager manager.Manager
}

func (d *driver) WatchEvents(ctx context.Context, types ...containers.EventType) (containers.EventsChannel, error) {
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

func convertContainerInfo(info v1.ContainerInfo) *containers.ContainerInfo {
	return &containers.ContainerInfo{
		Name:   info.Name,
		Labels: info.Labels,
	}
}

func (d *driver) GetContainers(ctx context.Context) ([]*containers.ContainerInfo, error) {
	q := &v1.ContainerInfoRequest{}
	rawContainers, err := d.manager.AllDockerContainers(q)
	if err != nil {
		return nil, err
	}

	result := make([]*containers.ContainerInfo, 0)
	for _, info := range rawContainers {
		result = append(result, convertContainerInfo(info))
	}

	return result, nil
}

func (d *driver) GetContainer(ctx context.Context, name string) (*containers.ContainerInfo, error) {
	r := &v1.ContainerInfoRequest{NumStats: 0}
	info, err := d.manager.GetContainerInfo(name, r)
	if err != nil {
		return nil, err
	}

	return convertContainerInfo(*info), nil
}

func (d *driver) GetContainerStats(ctx context.Context, container *containers.ContainerInfo) (containers.StatsChannel, error) {
	if !d.manager.Exists(container.Name) {
		return nil, errors.New("container does not exist")
	}

	return newStatsChannel(d.manager, container), nil
}

func (d *driver) CloseAllChannels(ctx context.Context) error {
	return nil
}
