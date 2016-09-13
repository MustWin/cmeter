package logger

import (
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
)

const NAME = "logger"

type Filter struct{}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx *pipeline.Context, m pipeline.Message) error {
	context.GetLogger(ctx).Debugf("processing %q message", m.Type())
	return nil
}

func New() pipeline.Filter {
	return &Filter{}
}
