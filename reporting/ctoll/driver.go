package ctoll

import (
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

	apiKey, _ := parameters["apikey"].(string)
	keyLabel, _ := parameters["key_label"].(string)
	return &Driver{
		keyLabel: keyLabel,
		client:   ctollclient.New(endpoint, apiKey, http.DefaultClient),
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

func calculateUsage(usage *containers.Usage) *v1.Usage {
	return &v1.Usage{
		TotalCPUPerc:   float64(usage.Cpu.Total),
		MemoryBytes:    usage.Memory.Bytes,
		DiskIOBytes:    sumOf(usage.Disk.PerDiskIo...),
		NetworkRxBytes: usage.Network.TotalRxBytes,
		NetworkTxBytes: usage.Network.TotalTxBytes,
	}
}

func calculateMachineUsage(u *containers.MachineUsage, m *containers.MachineInfo) *v1.MachineUsage {
	cores := float64(m.Cores)
	return &v1.MachineUsage{
		CPUShares:   (float64(u.Cpu.Total) / (1e+10 * cores)) * cores,
		MemoryBytes: u.Memory.Bytes,
	}
}

type Driver struct {
	keyLabel string
	client   *ctollclient.Client
}

func (d *Driver) Report(ctx context.Context, e *reporting.Event) (reporting.Receipt, error) {
	receiptData, err := d.sendEvent(e)
	if err != nil {
		err = fmt.Errorf("error sending event: %v", err)
	}

	return reporting.Receipt(string(receiptData)), err
}

func (d *Driver) apiKeyFromLabel(labels map[string]string) string {
	if d.keyLabel == "" {
		return ""
	}

	if v, ok := labels[d.keyLabel]; ok {
		return v
	}

	return ""
}

func (d *Driver) sendMeterStart(me *v1.MeterEvent, ch *containers.StateChange) ([]byte, error) {
	me.Type = v1.MeterEventTypeStart

	e := v1.StartMeterEvent{
		MeterEvent: me,
		Container:  convertContainerInfo(ch.Container),
		Allocated: &v1.BlockAlloc{
			CPUShares:   ch.Container.Reserved.Cpu,
			MemoryBytes: ch.Container.Reserved.Memory,
		},
	}

	key := d.apiKeyFromLabel(ch.Container.Labels)
	return []byte{}, d.client.MeterEvents().SendStartMeter(key, e)
}

func (d *Driver) sendMeterStop(me *v1.MeterEvent, ch *containers.StateChange) ([]byte, error) {
	me.Type = v1.MeterEventTypeStop

	e := v1.StopMeterEvent{MeterEvent: me}
	e.Container = convertContainerInfo(ch.Container)
	key := d.apiKeyFromLabel(ch.Container.Labels)
	return []byte{}, d.client.MeterEvents().SendStopMeter(key, e)
}

func (d *Driver) sendMeterSample(me *v1.MeterEvent, s *collector.Sample) ([]byte, error) {
	me.Type = v1.MeterEventTypeSample

	e := v1.SampleMeterEvent{
		MeterEvent: me,
		Usage:      calculateUsage(s.Usage),
		Container:  convertContainerInfo(s.Container),
	}

	key := d.apiKeyFromLabel(s.Container.Labels)
	return []byte{}, d.client.MeterEvents().SendUsageSample(key, e)
}

func (d *Driver) sendMeterMachineSample(me *v1.MeterEvent, s *collector.MachineSample) ([]byte, error) {
	me.Type = v1.MeterEventTypeMachineSample

	e := v1.MachineSampleMeterEvent{
		MeterEvent: me,
		Machine:    convertMachineInfo(s.Machine),
		Usage:      calculateMachineUsage(s.Usage, s.Machine),
	}

	key := d.apiKeyFromLabel(s.Machine.Labels)
	return []byte{}, d.client.MeterEvents().SendMachineUsageSample(key, e)
}

func (d *Driver) sendEvent(e *reporting.Event) ([]byte, error) {
	me := &v1.MeterEvent{
		MeterID:   e.MeterID,
		Timestamp: e.Timestamp,
	}

	switch e.Type {
	case reporting.EventStateChange:
		ch := e.Data.(*containers.StateChange)
		if ch.State == containers.StateRunning {
			return d.sendMeterStart(me, ch)
		} else {
			return d.sendMeterStop(me, ch)
		}

	case reporting.EventSample:
		return d.sendMeterSample(me, e.Data.(*collector.Sample))

	case reporting.EventMachineSample:
		return d.sendMeterMachineSample(me, e.Data.(*collector.MachineSample))
	}

	return []byte{}, fmt.Errorf("unsupported event type %q", e.Type)
}
