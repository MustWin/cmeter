package reporting

import (
	"time"

	"github.com/MustWin/cmeter/context"
)

type Receipt string

const EmptyReceipt Receipt = ""

type Event struct {
	MeterID   string      `json:"meter_id"`
	Type      string      `json:"event_type"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

var (
	EventSample        = "stat_sample"
	EventStateChange   = "state_change"
	EventMachineSample = "machine_stat_sample"
)

type Driver interface {
	Report(ctx context.Context, e *Event) (Receipt, error)
}

func Generate(ctx context.Context, eventType string, data interface{}) *Event {
	return &Event{
		Timestamp: time.Now().Unix(),
		MeterID:   context.GetInstanceID(ctx),
		Type:      eventType,
		Data:      data,
	}
}
