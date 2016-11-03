package embedded

import (
	"flag"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/cadvisor/cache/memory"
	cadvisorMetrics "github.com/google/cadvisor/container"
	"github.com/google/cadvisor/events"
	"github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/info/v2"
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
	defaultHousekeepingInterval = 5 * time.Second
	allowDynamicHousekeeping    = true

	sharesPerCPU = 1024.0
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

	machine, err := m.GetMachineInfo()
	if err != nil {
		return nil, err
	}

	cpuLimitLabel, _ := parameters["cpu_limit_label"].(string)

	d := &driver{
		cpuLimitLabel: cpuLimitLabel,
		machine:       convertMachineInfo(machine),
		manager:       m,
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
			// TODO: can't log error here but marking as a *maybe* and might find a better approach for this
			//log.Errorf("Expected cAdvisor flag %q not found", name)
		}
	}
}

type driver struct {
	cpuLimitLabel string
	manager       manager.Manager
	machine       *containers.MachineInfo
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

func parseImageData(image string) (string, string) {
	parts := strings.Split(image, ":")
	if len(parts) < 2 {
		return image, ""
	} else if parts[1] == "" {
		parts[1] = "latest"
	}

	return parts[0], parts[1]
}

func maxCpuLimit(shares float64, cores int) float64 {
	return shares / sharesPerCPU
}

func maxCpuLimitOverride(labels map[string]string, limitLabel string) float64 {
	if limitLabel == "" {
		return 0
	}

	limitStr, ok := labels[limitLabel]
	if !ok || limitStr == "" {
		return 0
	}

	f, err := strconv.ParseFloat(limitStr, 64)
	if err != nil {
		return 0
	}

	return f
}

func convertMachineInfo(info *v1.MachineInfo) *containers.MachineInfo {
	return &containers.MachineInfo{
		SystemUuid:      info.SystemUUID,
		Cores:           info.NumCores,
		MemoryBytes:     info.MemoryCapacity,
		CpuFrequencyKhz: info.CpuFrequency,
	}
}

func convertContainerInfo(info v1.ContainerInfo, machine *containers.MachineInfo, cpuLimitLabel string) *containers.ContainerInfo {
	imageName, imageTag := parseImageData(info.Spec.Image)
	cpuLimit := maxCpuLimit(float64(info.Spec.Cpu.Limit), machine.Cores)
	if cpuLimitLabel != "" {
		cpuLimit = maxCpuLimitOverride(info.Labels, cpuLimitLabel)
	}

	return &containers.ContainerInfo{
		Name:      info.Name,
		ImageName: imageName,
		ImageTag:  imageTag,
		Labels:    info.Labels,
		Machine:   machine,
		Reserved: &containers.ReservedResources{
			Cpu:    cpuLimit,
			Memory: info.Spec.Memory.Limit,
		},
	}
}

func convertContainerSpec(name string, spec v2.ContainerSpec, machine *containers.MachineInfo, cpuLimitLabel string) *containers.ContainerInfo {
	imageName, imageTag := parseImageData(spec.Image)
	cpuLimit := maxCpuLimit(float64(spec.Cpu.Limit), machine.Cores)
	if cpuLimitLabel != "" {
		cpuLimit = maxCpuLimitOverride(spec.Labels, cpuLimitLabel)
	}

	return &containers.ContainerInfo{
		Name:      name,
		ImageName: imageName,
		ImageTag:  imageTag,
		Labels:    spec.Labels,
		Machine:   machine,
		Reserved: &containers.ReservedResources{
			Cpu:    cpuLimit,
			Memory: spec.Memory.Limit,
		},
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
		result = append(result, convertContainerInfo(info, d.machine, d.cpuLimitLabel))
	}

	return result, nil
}

func (d *driver) GetContainer(ctx context.Context, name string) (*containers.ContainerInfo, error) {
	specMap, err := d.manager.GetContainerSpec(name, v2.RequestOptions{
		IdType:    "name",
		Count:     0,
		Recursive: false,
	})

	if err != nil {
		if strings.Contains(err.Error(), "unable to find data for container") {
			return nil, containers.ErrContainerNotFound
		}

		return nil, err
	}

	return convertContainerSpec(name, specMap[name], d.machine, d.cpuLimitLabel), nil
}

func (d *driver) GetContainerStats(ctx context.Context, name string) (containers.StatsChannel, error) {
	container, err := d.GetContainer(ctx, name)
	if err != nil {
		return nil, err
	}

	return newStatsChannel(d.manager, container), nil
}

func (d *driver) CloseAllChannels(ctx context.Context) error {
	// TODO: determine need and complete
	return nil
}
