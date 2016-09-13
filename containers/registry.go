package containers

import (
	"github.com/MustWin/cmeter/context"
)

type Registry struct {
	containers map[string]*ContainerInfo
}

func NewRegistry() *Registry {
	return &Registry{
		containers: make(map[string]*ContainerInfo),
	}
}

func (registry *Registry) Get(containerName string) (*ContainerInfo, bool) {
	info, ok := registry.containers[containerName]
	return info, ok
}

func (registry *Registry) IsRegistered(containerName string) bool {
	_, ok := registry.containers[containerName]
	return ok
}

func (registry *Registry) Register(ctx context.Context, info *ContainerInfo) error {
	log := context.GetLoggerWithField(ctx, "container.name", info.Name)
	if registry.IsRegistered(info.Name) {
		log.Warnf("container already registered with name %q, ignoring", info.Name)
		return nil
	}

	registry.containers[info.Name] = info
	log.Infof("registered container %q", info.Name)
	return nil
}
