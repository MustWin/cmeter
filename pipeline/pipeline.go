package pipeline

import (
	"fmt"

	"github.com/MustWin/cmeter/context"
)

type Message interface {
	Type() string
	Body() interface{}
}

type Pipeline interface {
	Send(ctx context.Context, m Message)
}

type Filter interface {
	Name() string
	HandleMessage(ctx context.Context, m Message) error
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
	var err error
	ctx, err = StartProcessing(pipe, ctx, m)
	if err != nil {
		context.GetLogger(ctx).Errorf("error start processing: %v", err)
		return
	}

	processingID, _ := GetProcessingID(ctx)
	ctx = context.WithLogger(ctx, context.GetLoggerWithFields(ctx, map[interface{}]interface{}{
		"send.id":      processingID,
		"message.type": m.Type(),
	}))

	for _, filter := range pipe.filters {
		err := filter.HandleMessage(ctx, m)
		if err != nil {
			context.GetLoggerWithField(ctx, "filter.name", filter.Name()).Errorf("filter error processing message: %v", err)
			if err := StopProcessing(ctx); err != nil {
				context.GetLogger(ctx).Errorf("error stop processing on context: %v", err)
			}

			break
		}

		if !IsProcessing(ctx) {
			break
		}

		m, _ = GetMessage(ctx)
	}
}

func New(filters ...Filter) Pipeline {
	return &simplePipe{filters: filters}
}
