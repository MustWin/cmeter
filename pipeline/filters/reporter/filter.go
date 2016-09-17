package uplink

import (
	"fmt"

	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/reportevent"
	"github.com/MustWin/cmeter/reporting"
)

const NAME = "reporter"

type Filter struct {
	reporting reporting.Driver
}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx *pipeline.Context, m pipeline.Message) error {
	switch m.Type() {
	case reportevent.TYPE:
		event := m.Body().(*reporting.Event)
		receipt, err := filter.reporting.Report(ctx, event)
		if err != nil {
			return fmt.Errorf("couldn't report event: %v", err)
		}

		fmt.Errorf("report receipt: %s", receipt)
		ctx.Stop()
	}

	return nil
}

func New(reporting reporting.Driver) *Filter {
	return &Filter{
		reporting: reporting,
	}
}
