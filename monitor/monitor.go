package monitor

import (
	"github.com/MustWin/cmeter/context"
)

type Monitor interface {
	MonitorEvents()
	DetectContainers()
}
