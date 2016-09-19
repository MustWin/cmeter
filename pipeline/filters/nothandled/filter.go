package nothandled

import (
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
)

const NAME = "not_handled"

type Filter struct{}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx *pipeline.Context, m pipeline.Message) error {
	context.GetLogger(ctx).Warn("message not handled")
	ctx.Stop()
	return nil
}

func New() pipeline.Filter {
	return &Filter{}
}
