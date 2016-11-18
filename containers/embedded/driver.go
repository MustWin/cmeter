package embedded

import (
	"flag"
	"math"
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

var (
	parseOnce sync.Once
	terabyte  uint64
)

const (
	statsCacheDuration          = 2 * time.Minute
	maxHousekeepingInterval     = 15 * time.Second
	defaultHousekeepingInterval = 5 * time.Second
	allowDynamicHousekeeping    = true

	sharesPerCPU = 1024.0

	rootContainerName = "/"
)

func init() {
	terabyte = uint64(math.Pow(1024, 4))

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

	// override before we instantiate the manager
	allowedEnvs := getEnvWhiteList(parameters)
	f := flag.Lookup("docker_env_metadata_whitelist")
	f.Value.Set(strings.Join(allowedEnvs, ","))

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

	if err = m.Start(); err != nil {
		return nil, err
	}

	rootMap, err := m.GetContainerSpec(rootContainerName, v2.RequestOptions{
		IdType:    v2.TypeName,
		Count:     0,
		Recursive: false,
	})

	if err != nil {
		return nil, err
	}

	d := &driver{
		cpuLimitLabel: cpuLimitLabel,
		machine:       convertMachineInfo(machine, rootMap[rootContainerName]),
		manager:       m,
	}

	return d, nil
}

func getEnvWhiteList(parameters map[string]interface{}) []string {
	if delimited, ok := parameters["envs"].(string); ok {
		return strings.Split(delimited, ",")
	}

	if raw, ok := parameters["envs"].([]interface{}); ok {
		envs := make([]string, 0)
		for _, e := range raw {
			if s, ok := e.(string); ok && s != "" {
				envs = append(envs, s)
			}
		}

		return envs
	}

	return nil
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

func maxCpuLimitOverride(limit float64, labels map[string]string, limitLabel string) float64 {
	if limitLabel == "" {
		return limit
	}

	limitStr, ok := labels[limitLabel]
	if !ok || limitStr == "" {
		return limit
	}

	f, err := strconv.ParseFloat(limitStr, 64)
	if err != nil {
		return limit
	}

	return f
}

func normalizeMemoryLimit(limit uint64) uint64 {
	// unreasonably high amount of memory allowance, assume unbounded
	if limit >= terabyte {
		return 0
	}

	return limit
}

func convertMachineInfo(info *v1.MachineInfo, rootSpec v2.ContainerSpec) *containers.MachineInfo {
	name := ""
	if info.InstanceID != v1.UnNamedInstance {
		name = string(info.InstanceID)
	}

	return &containers.MachineInfo{
		SystemUuid:      info.SystemUUID,
		Cores:           info.NumCores,
		MemoryBytes:     info.MemoryCapacity,
		CpuFrequencyKhz: info.CpuFrequency,
		Labels:          rootSpec.Labels,
		Name:            name,
	}
}

func convertContainerInfo(info v1.ContainerInfo, machine *containers.MachineInfo, cpuLimitLabel string) *containers.ContainerInfo {
	imageName, imageTag := parseImageData(info.Spec.Image)
	cpuLimit := maxCpuLimit(float64(info.Spec.Cpu.Limit), machine.Cores)
	if cpuLimitLabel != "" {
		cpuLimit = maxCpuLimitOverride(cpuLimit, info.Labels, cpuLimitLabel)
	}

	return &containers.ContainerInfo{
		Name:      info.Name,
		ImageName: imageName,
		ImageTag:  imageTag,
		Labels:    info.Labels,
		Machine:   machine,
		Envs:      info.Spec.Envs,
		Reserved: &containers.ReservedResources{
			Cpu:    cpuLimit,
			Memory: normalizeMemoryLimit(info.Spec.Memory.Limit),
		},
	}
}

func convertContainerSpec(name string, spec v2.ContainerSpec, machine *containers.MachineInfo, cpuLimitLabel string) *containers.ContainerInfo {
	imageName, imageTag := parseImageData(spec.Image)
	cpuLimit := maxCpuLimit(float64(spec.Cpu.Limit), machine.Cores)
	if cpuLimitLabel != "" {
		cpuLimit = maxCpuLimitOverride(cpuLimit, spec.Labels, cpuLimitLabel)
	}

	return &containers.ContainerInfo{
		Name:      name,
		ImageName: imageName,
		ImageTag:  imageTag,
		Labels:    spec.Labels,
		Machine:   machine,
		Envs:      spec.Envs,
		Reserved: &containers.ReservedResources{
			Cpu:    cpuLimit,
			Memory: normalizeMemoryLimit(spec.Memory.Limit),
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
		IdType:    v2.TypeName,
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

func (d *driver) GetContainerUsage(ctx context.Context, name string) (containers.UsageChannel, error) {
	container, err := d.GetContainer(ctx, name)
	if err != nil {
		return nil, err
	}

	return newUsageChannel(d.manager, container), nil
}

func (d *driver) CloseAllChannels(ctx context.Context) error {
	// TODO: determine need and complete
	return nil
}

func (d *driver) GetMachineUsage(ctx context.Context) (containers.MachineUsageFeed, error) {
	root, err := d.GetContainer(ctx, rootContainerName)
	if err != nil {
		return nil, err
	}

	return newMachineUsageFeed(d.manager, d.machine, root), nil
}
