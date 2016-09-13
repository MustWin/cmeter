package embedded

import (
	"github.com/MustWin/cmeter/configuration"
)

type Tracker struct {
	containers      map[string]*ContainerInfo
	keyLabel        string
	serviceKeyLabel string
}

func NewTracker(config *configuration.TrackerConfig) *ContainerTracker {
	return &Tracker{
		containers:      make(map[string]*ContainerInfo),
		keyLabel:        config.KeyLabel,
		serviceKeyLabel: config.ServiceKeyLabel,
	}
}
