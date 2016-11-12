package containers

import (
	"errors"
	"strings"
	"sync"

	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/context"
)

var ErrNotTrackable = errors.New("the registry cannot track this container")

type Registry struct {
	mutex      sync.Mutex
	markers    configuration.Marker
	containers map[string]*ContainerInfo
}

func NewRegistry(markers configuration.Marker) *Registry {
	markers.Env = strings.ToLower(markers.Env)
	return &Registry{
		containers: make(map[string]*ContainerInfo),
		markers:    markers,
	}
}

func (registry *Registry) Get(containerName string) (*ContainerInfo, bool) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	info, ok := registry.containers[containerName]
	return info, ok
}

func (registry *Registry) List() []*ContainerInfo {
	results := make([]*ContainerInfo, len(registry.containers))
	i := 0
	for _, c := range registry.containers {
		results[i] = c
		i++
	}

	return results
}

func (registry *Registry) IsRegistered(containerName string) bool {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	_, ok := registry.containers[containerName]
	return ok
}

func (registry *Registry) Register(ctx context.Context, info *ContainerInfo) error {
	if !registry.IsTrackable(info) {
		return ErrNotTrackable
	}

	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	log := context.GetLoggerWithField(ctx, "container.name", info.Name)
	if _, ok := registry.containers[info.Name]; ok {
		log.Warnf("container name already registered, ignoring", info.Name)
		return nil
	}

	registry.containers[info.Name] = info
	log.Info("container registered")
	return nil
}

func (registry *Registry) Drop(ctx context.Context, containerName string) error {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	log := context.GetLoggerWithField(ctx, "container.name", containerName)
	if _, ok := registry.containers[containerName]; !ok {
		log.Warnf("container name not registered, ignoring", containerName)
		return nil
	}

	delete(registry.containers, containerName)
	log.Info("container dropped")
	return nil
}

func (registry *Registry) IsTrackable(info *ContainerInfo) bool {
	if registry.markers.Label != "" {
		if _, ok := info.Labels[registry.markers.Label]; ok {
			return true
		}
	}

	if registry.markers.Env != "" {
		if _, ok := info.Envs[registry.markers.Env]; ok {
			return true
		}
	}

	return false
}
