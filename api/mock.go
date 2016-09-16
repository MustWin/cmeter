package api

import (
	"github.com/MustWin/cmeter/context"
)

type mockClient struct{}

func (c *mockClient) Send(ctx context.Context, e *Event) error {
	context.GetLogger(ctx).Infof("transmit event: %+#v", e)
	return nil
}

func NewMockClient() Client {
	return &mockClient{}
}
