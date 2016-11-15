package embedded

import (
	"errors"
	"sync"

	"github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/manager"

	"github.com/MustWin/cmeter/containers"
)

type usageChannel struct {
	startFetch sync.Once
	manager    manager.Manager
	container  *containers.ContainerInfo
	ch         chan *containers.Usage
	doneCh     chan bool
	closed     bool
	last       *v1.ContainerStats
}

func (ch *usageChannel) Container() *containers.ContainerInfo {
	return ch.container
}

func (ch *usageChannel) GetChannel() <-chan *containers.Usage {
	ch.startFetch.Do(func() {
		go ch.startChannel()
	})

	return ch.ch
}

func (ch *usageChannel) startChannel() {
	defer close(ch.ch)
	for {
		select {
		case done, ok := <-ch.doneCh:
			if done || !ok {
				return
			}

		default:
			ci, err := ch.manager.GetContainerInfo(ch.container.Name, &v1.ContainerInfoRequest{NumStats: 1})
			if err == nil && ci != nil && len(ci.Stats) > 0 {
				stats := ci.Stats[0]
				if ch.last == nil {
					ch.last = stats
				}

				cs := convertContainerInfoToStats(ch.last, stats)
				ch.last = stats
				ch.ch <- cs
			}
		}
	}
}

func (ch *usageChannel) Close() error {
	select {
	case _, ok := <-ch.doneCh:
		if !ok {
			return errors.New("stats channel already closed")
		}
	}

	ch.doneCh <- true
	close(ch.doneCh)
	return nil
}

type machineUsageFeed struct {
	machine *containers.MachineInfo
	root    *containers.ContainerInfo
	manager manager.Manager
	last    *v1.ContainerStats
}

func (ch *machineUsageFeed) Next() *containers.MachineUsage {
	ci, err := ch.manager.GetContainerInfo(ch.root.Name, &v1.ContainerInfoRequest{NumStats: 1})
	if err == nil && ci != nil && len(ci.Stats) > 0 {
		stats := ci.Stats[0]
		if ch.last == nil {
			ch.last = stats
		}

		ms := getMachineUsage(ch.last, stats)
		ch.last = stats
		return ms
	}

	return nil
}

func (ch *machineUsageFeed) Machine() *containers.MachineInfo {
	return ch.machine
}

func newMachineUsageFeed(manager manager.Manager, machine *containers.MachineInfo, root *containers.ContainerInfo) *machineUsageFeed {
	f := &machineUsageFeed{
		manager: manager,
		root:    root,
		machine: machine,
	}

	// prime it
	f.Next()
	return f
}

func newUsageChannel(manager manager.Manager, container *containers.ContainerInfo) *usageChannel {
	ch := &usageChannel{
		manager:   manager,
		container: container,
		closed:    false,
		ch:        make(chan *containers.Usage),
		doneCh:    make(chan bool),
	}

	// prime it
	<-ch.GetChannel()
	return ch
}

func getMachineUsage(last, stats *v1.ContainerStats) *containers.MachineUsage {
	cu := convertContainerInfoToStats(last, stats)
	return &containers.MachineUsage{
		Cpu:    cu.Cpu,
		Memory: cu.Memory,
	}
}

func convertContainerInfoToStats(last, stats *v1.ContainerStats) *containers.Usage {
	cpu := &containers.CpuUsage{
		Total:   int64(stats.Cpu.Usage.Total - last.Cpu.Usage.Total),
		PerCore: make([]int64, len(stats.Cpu.Usage.PerCpu)),
	}

	for i, coreNs := range stats.Cpu.Usage.PerCpu {
		cpu.PerCore[i] = int64(coreNs - last.Cpu.Usage.PerCpu[i])
	}

	disk := &containers.DiskUsage{
		PerDiskIo: make([]uint64, 0),
	}

	for _, d := range stats.DiskIo.IoServiceBytes {
		total := uint64(0)
		for _, bytesTransferred := range d.Stats {
			total += bytesTransferred
		}

		disk.PerDiskIo = append(disk.PerDiskIo, total)
	}

	net := &containers.NetworkUsage{
		TotalRxBytes: stats.Network.RxBytes,
		TotalTxBytes: stats.Network.TxBytes,
		Interfaces:   make([]*containers.InterfaceUsage, 0),
	}

	for _, nic := range stats.Network.Interfaces {
		net.Interfaces = append(net.Interfaces, &containers.InterfaceUsage{
			Name:    nic.Name,
			RxBytes: nic.RxBytes,
			TxBytes: nic.TxBytes,
		})
	}

	return &containers.Usage{
		Cpu:     cpu,
		Memory:  &containers.MemoryUsage{Bytes: stats.Memory.Usage},
		Disk:    disk,
		Network: net,
	}
}
