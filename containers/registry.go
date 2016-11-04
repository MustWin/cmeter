package containers

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/MustWin/cmeter/context"
)

var ErrNotTrackable = errors.New("the registry cannot track this container")

type Registry struct {
	mutex          sync.Mutex
	trackingLabel  string
	trackingEnvKey string
	containers     map[string]*ContainerInfo
}

func NewRegistry(trackingLabel string, trackingEnvKey string) *Registry {
	return &Registry{
		containers:     make(map[string]*ContainerInfo),
		trackingLabel:  trackingLabel,
		trackingEnvKey: strings.ToLower(trackingEnvKey),
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
	if registry.trackingLabel != "" {
		if _, ok := info.Labels[registry.trackingLabel]; ok {
			return true
		}
	}

	if registry.trackingEnvKey != "" {
		fmt.Printf("key/env: %s  |  %#+v\n", registry.trackingEnvKey, info.Envs)
		if _, ok := info.Envs[registry.trackingEnvKey]; ok {
			return true
		}
	}

	return false
}
