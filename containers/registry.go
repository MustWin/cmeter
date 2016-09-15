package containers

import (
	"sync"

	"github.com/MustWin/cmeter/context"
)

type Registry struct {
	mutex      sync.Mutex
	containers map[string]*ContainerInfo
}

func NewRegistry() *Registry {
	return &Registry{
		containers: make(map[string]*ContainerInfo),
	}
}

func (registry *Registry) Get(containerName string) (*ContainerInfo, bool) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	info, ok := registry.containers[containerName]
	return info, ok
}

func (registry *Registry) IsRegistered(containerName string) bool {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	_, ok := registry.containers[containerName]
	return ok
}

func (registry *Registry) Register(ctx context.Context, info *ContainerInfo) error {
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
