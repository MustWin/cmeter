package mock

import (
	"fmt"

	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/reporting"
	"github.com/MustWin/cmeter/reporting/factory"
)

type driverFactory struct{}

func (factory *driverFactory) Create(_ map[string]interface{}) (reporting.Driver, error) {
	return &Driver{}, nil
}

func init() {
	factory.Register("mock", &driverFactory{})
}

type Driver struct {
	Counter uint64
}

func (d *Driver) Report(ctx context.Context, e *reporting.Event) (reporting.Receipt, error) {
	d.Counter++
	receipt := reporting.Receipt(fmt.Sprintf("event#%d", d.Counter))
	context.GetLogger(ctx).Info("report@%d: %+#v", d.Counter, e)
	return reporting.Receipt(receipt), nil
}
