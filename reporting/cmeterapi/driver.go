package cmeterapi

import (
	"errors"

	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/reporting"
	"github.com/MustWin/cmeter/reporting/factory"
)

var ErrInvalidRemoteAddr = errors.New("invalid remote address")

type driverFactory struct{}

func (factory *driverFactory) Create(parameters map[string]interface{}) (reporting.Driver, error) {
	remoteAddr, ok := parameters["addr"].(string)
	if !ok || remoteAddr == "" {
		return nil, ErrInvalidRemoteAddr
	}

	return &Driver{
		RemoteAddr: remoteAddr,
	}, nil
}

func init() {
	factory.Register("cmeterapi", &driverFactory{})
}

type Driver struct {
	RemoteAddr string
}

func (d *Driver) Report(ctx context.Context, e *reporting.Event) (reporting.Receipt, error) {
	// TODO: implementation
	return reporting.Receipt(""), nil
}
