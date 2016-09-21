package reporter

import (
	"fmt"

	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/messages/sendreport"
	"github.com/MustWin/cmeter/reporting"
)

const (
	NAME = "reporter"
)

type Filter struct {
	driver reporting.Driver
}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx context.Context, m pipeline.Message) error {
	if m.Type() == sendreport.TYPE {
		receipt, err := filter.driver.Report(ctx, m.Body().(*reporting.Event))
		if err != nil {
			return fmt.Errorf("error performing report: %v", err)
		}

		context.GetLoggerWithField(ctx, "report.receipt", receipt).Info("sent report")
		return pipeline.StopProcessing(ctx)
	}

	return nil
}

var _ pipeline.Filter = &Filter{}

func New(driver reporting.Driver) *Filter {
	return &Filter{
		driver: driver,
	}
}
