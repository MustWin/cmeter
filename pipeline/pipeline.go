package pipeline

import (
	"fmt"

	"github.com/MustWin/cmeter/context"
)

type Context struct {
	context.Context
	stopped  bool
	Pipeline Pipeline
}

func (ctx *Context) Stop() {
	ctx.stopped = true
}

func (ctx *Context) Stopped() bool {
	return ctx.stopped
}

type Message interface {
	ID() string
	Type() string
	Body() interface{}
}

type Pipeline interface {
	Send(ctx context.Context, m Message) error
}

type Filter interface {
	Name() string
	HandleMessage(ctx *Context, m Message) error
}

type FilterError struct {
	FilterName string
	Enclosed   error
}

func (err FilterError) Error() string {
	return fmt.Sprintf("FilterError: %s: %v", err.FilterName, err.Enclosed)
}

type simplePipe struct {
	filters []Filter
}

func (pipe *simplePipe) Send(ctx context.Context, m Message) error {
	ctx = context.WithLogger(ctx, context.GetLoggerWithFields(ctx, map[interface{}]interface{}{
		"message.id":   m.ID(),
		"message.type": m.Type(),
	}))

	pctx := &Context{
		Context: ctx,
		stopped: false,
	}

	for _, filter := range pipe.filters {
		err := filter.HandleMessage(pctx, m)
		if err != nil {
			return FilterError{
				FilterName: filter.Name(),
				Enclosed:   err,
			}
		}

		if pctx.Stopped() {
			break
		}
	}

	return nil
}

func New(filters ...Filter) Pipeline {
	return &simplePipe{filters: filters}
}
