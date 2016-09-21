package reportgen

import (
	"fmt"
	"time"

	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/containersample"
	"github.com/MustWin/cmeter/pipeline/messages/sendreport"
	"github.com/MustWin/cmeter/pipeline/messages/statechange"
	"github.com/MustWin/cmeter/reporting"
)

const (
	NAME = "report_generator"

	EVENT_SAMPLE       = "stat_sample"
	EVENT_STATE_CHANGE = "state_change"
)

type Filter struct{}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx context.Context, m pipeline.Message) error {
	var r *reporting.Event
	switch m.Type() {
	case containersample.TYPE:
		r = &reporting.Event{}
		r.Type = EVENT_SAMPLE
		r.Data = m.Body()

	case statechange.TYPE:
		r = &reporting.Event{}
		r.Type = EVENT_STATE_CHANGE
		r.Data = m.Body()
	}

	if r != nil {
		r.MeterID = context.GetInstanceID(ctx)
		r.Timestamp = time.Now().Unix()
		if err := pipeline.SetMessage(ctx, sendreport.NewMessage(r)); err != nil {
			return fmt.Errorf("couldn't set downstream message: %v", err)
		}
	}

	return nil
}

var _ pipeline.Filter = &Filter{}

func New() *Filter {
	return &Filter{}
}
