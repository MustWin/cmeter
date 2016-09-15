package pipeline

import (
	"fmt"

	"github.com/MustWin/cmeter/context"
)

type Context struct {
	context.Context
	stopped         bool
	messageOverride Message
	Pipeline        Pipeline
}

func (ctx *Context) SetMessage(m Message) {
	ctx.messageOverride = m
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
	Send(ctx context.Context, m Message)
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

func (pipe *simplePipe) Send(ctx context.Context, m Message) {
	ctx = context.WithLogger(ctx, context.GetLoggerWithFields(ctx, map[interface{}]interface{}{
		"message.id":   m.ID(),
		"message.type": m.Type(),
	}))

	pctx := &Context{
		Pipeline: pipe,
		Context:  ctx,
		stopped:  false,
	}

	for _, filter := range pipe.filters {
		err := filter.HandleMessage(pctx, m)
		if err != nil {
			// TODO: send error message
			context.GetLoggerWithField(ctx, "filter.name", filter.Name()).Errorf("filter error processing message: %v", err)
			break
		}

		if pctx.messageOverride != nil {
			m = pctx.messageOverride
			pctx.messageOverride = nil
			pctx.Context = context.WithLogger(ctx, context.GetLoggerWithFields(ctx, map[interface{}]interface{}{
				"message.id":   m.ID(),
				"message.type": m.Type(),
			}))
		}

		if pctx.Stopped() {
			break
		}
	}
}

func New(filters ...Filter) Pipeline {
	return &simplePipe{filters: filters}
}
