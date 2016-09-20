package pipeline

import (
	"fmt"

	"github.com/MustWin/cmeter/context"
)

var ErrCtxNoPipeline = errors.New("no pipeline associated with this context")

func getPipelineContext(ctx context.Context) (*context, error) {
	if pctx, ok := ctx.Value("pipeline.ctx").(*context); ok {
		return pctx
	}

	return ErrCtxNoPipeline
}

func StopProcessing(ctx context.Context) error {
	pctx, err := getPipelineContext(ctx)
	if err != nil {
		return err
	}

	pctx.stopped = true
	return nil
}

type context struct {
	context.Context
	stopped         bool
	messageOverride Message
	pipeline        Pipeline
}

func (ctx *context) Value(key interface{}) interface{} {
	switch key {
	case "pipeline":
		return ctx.Pipeline

	case "pipeline.ctx":
		return ctx
	}

	if value, ok := ctx.Value(key).(string); ok {
		return value
	}

	return ctx.Context.Value(key)
}

func (ctx *context) SetMessage(m Message) {
	ctx.messageOverride = m
}

func (ctx *context) Stopped() bool {
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
