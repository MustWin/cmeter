package api

import (
	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/context"
)

type Client interface {
	Send(ctx context.Context, e *Event) error
}

type Event struct {
	MeterID    string      `json:"meter_id"`
	Type       string      `json:"event_type"`
	ServiceKey string      `json:"service_key"`
	Timestamp  int64       `json:"timestamp"`
	Data       interface{} `json:"data"`
}

type client struct {
	remoteAddr string
}

func (c *client) Send(ctx context.Context, e *Event) error {
	return nil
}

func NewClient(config configuration.ApiConfig) Client {
	return &client{
		remoteAddr: config.Addr,
	}
}
