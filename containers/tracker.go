package containers

import (
	"github.com/MustWin/cmeter/configuration"
)

type Tracker struct {
	containers      map[string]*ContainerInfo
	trackingLabel   string
	serviceKeyLabel string
}

func NewTracker(config *configuration.TrackerConfig) *Tracker {
	return &Tracker{
		containers:      make(map[string]*ContainerInfo),
		trackingLabel:   config.TrackingLabel,
		serviceKeyLabel: config.ServiceKeyLabel,
	}
}
