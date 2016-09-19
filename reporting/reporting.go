package reporting

import (
	"github.com/MustWin/cmeter/context"
)

type Receipt string

type Event struct {
	MeterID   string      `json:"meter_id"`
	Type      string      `json:"event_type"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type Driver interface {
	Report(ctx context.Context, e *Event) (Receipt, error)
}
