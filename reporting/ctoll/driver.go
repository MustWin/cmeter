package ctoll

import (
	"errors"
	"fmt"
	"net/http"

	ctollclient "github.com/MustWin/ctoll/ctoll/api/client"
	"github.com/MustWin/ctoll/ctoll/api/v1"

	"github.com/MustWin/cmeter/collector"
	"github.com/MustWin/cmeter/containers"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/reporting"
	"github.com/MustWin/cmeter/reporting/factory"
)

type driverFactory struct{}

func (factory *driverFactory) Create(parameters map[string]interface{}) (reporting.Driver, error) {
	endpoint, ok := parameters["endpoint"].(string)
	if !ok || endpoint == "" {
		// TODO: default to whatever public url ctoll eventually gets
		endpoint = "http://localhost:9180"
	}

	apiKey, ok := parameters["apikey"].(string)
	if !ok || apiKey == "" {
		return nil, errors.New("cToll api key missing or invalid")
	}

	return &Driver{
		client: ctollclient.New(endpoint, apiKey, http.DefaultClient),
	}, nil
}

func init() {
	factory.Register("ctoll", &driverFactory{})
}

func convertMachineInfo(m *containers.MachineInfo) *v1.MachineInfo {
	return &v1.MachineInfo{
		SystemUuid:      m.SystemUuid,
		Cores:           m.Cores,
		CpuFrequencyKhz: m.CpuFrequencyKhz,
		MemoryBytes:     m.MemoryBytes,
	}
}

func convertContainerInfo(ci *containers.ContainerInfo) *v1.ContainerInfo {
	return &v1.ContainerInfo{
		ImageName: ci.ImageName,
		ImageTag:  ci.ImageTag,
		Name:      ci.Name,
		Labels:    ci.Labels,
		Machine:   convertMachineInfo(ci.Machine),
	}
}

func sumOf(values ...uint64) uint64 {
	sum := uint64(0)
	for _, v := range values {
		sum += uint64(v)
	}

	return sum
}

func calculateUsage(stats *containers.Stats) *v1.Usage {
	return &v1.Usage{
		TotalCPUPerc:   stats.Cpu.TotalUsagePerc,
		MemoryBytes:    stats.Memory.UsageBytes,
		DiskIOBytes:    sumOf(stats.Disk.PerDiskIoBytes...),
		NetworkRxBytes: stats.Network.TotalRxBytes,
		NetworkTxBytes: stats.Network.TotalTxBytes,
	}
}

type Driver struct {
	client *ctollclient.Client
}

func (d *Driver) Report(ctx context.Context, e *reporting.Event) (reporting.Receipt, error) {
	receiptData, err := d.sendEvent(e)
	if err != nil {
		err = fmt.Errorf("error sending event: %v", err)
	}

	return reporting.Receipt(string(receiptData)), err
}

func (d *Driver) sendMeterStart(me *v1.MeterEvent, ch *containers.StateChange) ([]byte, error) {
	me.Type = v1.MeterEventTypeStart

	e := v1.StartMeterEvent{
		MeterEvent: me,
		//Allocated: ,
	}

	return []byte{}, d.client.MeterEvents().SendStartMeter(e)
}

func (d *Driver) sendMeterStop(me *v1.MeterEvent, ch *containers.StateChange) ([]byte, error) {
	me.Type = v1.MeterEventTypeStop
	me.Container = convertContainerInfo(ch.Container)

	e := v1.StopMeterEvent{MeterEvent: me}
	return []byte{}, d.client.MeterEvents().SendStopMeter(e)
}

func (d *Driver) sendMeterSample(me *v1.MeterEvent, s *collector.Sample) ([]byte, error) {
	me.Type = v1.MeterEventTypeSample
	me.Container = convertContainerInfo(s.Container)

	e := v1.SampleMeterEvent{
		MeterEvent: me,
		Usage:      calculateUsage(s.Stats),
	}

	return []byte{}, d.client.MeterEvents().SendUsageSample(e)
}

func (d *Driver) sendEvent(e *reporting.Event) ([]byte, error) {
	me := &v1.MeterEvent{
		MeterID:   e.MeterID,
		Timestamp: e.Timestamp,
	}

	switch e.Type {
	case reporting.EventStateChange:
		ch := e.Data.(*containers.StateChange)
		me.Container = convertContainerInfo(ch.Container)
		if ch.State == containers.StateRunning {
			return d.sendMeterStart(me, ch)
		} else {
			return d.sendMeterStop(me, ch)
		}

	case reporting.EventSample:
		return d.sendMeterSample(me, e.Data.(*collector.Sample))
	}

	return []byte{}, fmt.Errorf("unsupported event type %q", e.Type)
}
