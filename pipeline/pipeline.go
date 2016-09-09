package pipeline

import (
	"fmt"

	"github.com/MustWin/cmeter/context"
)

type Context struct {
	context.Context
	Stopped bool
}

func (ctx *Context) Stop() {
	ctx.Stopped = true
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
		Stopped: false,
	}

	for _, filter := range pipe.filters {
		err := filter.HandleMessage(pctx, m)
		if err != nil {
			return FilterError{
				FilterName: filter.Name(),
				Enclosed:   err,
			}
		}

		if pctx.Stopped {
			break
		}
	}

	return nil
}

func New(filters ...Filter) Pipeline {
	return &simplePipe{filters: filters}
}
