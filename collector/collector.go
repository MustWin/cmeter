package collector

import (
	"github.com/MustWin/cmeter/context"
)

type Driver interface {
	GatherMetrics(ctx context.Context, pipeline pipeline.Pipeline)
}
