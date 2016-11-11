package embedded

import (
	"errors"
	"sync"

	"github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/manager"

	"github.com/MustWin/cmeter/containers"
)

type statsChannel struct {
	startFetch sync.Once
	manager    manager.Manager
	container  *containers.ContainerInfo
	ch         chan *containers.Stats
	doneCh     chan bool
	closed     bool
}

func (ch *statsChannel) Container() *containers.ContainerInfo {
	return ch.container
}

func (ch *statsChannel) GetChannel() <-chan *containers.Stats {
	ch.startFetch.Do(func() {
		go ch.startChannel()
	})

	return ch.ch
}

func (ch *statsChannel) startChannel() {
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
				ch.ch <- convertContainerInfoToStats(ci.Stats[0])
			}
		}
	}
}

func (ch *statsChannel) Close() error {
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

type machineStatsFeed struct {
	machine *containers.MachineInfo
	root    *containers.ContainerInfo
	manager manager.Manager
}

func (ch *machineStatsFeed) Next() *containers.MachineStats {
	ci, err := ch.manager.GetContainerInfo(ch.root.Name, &v1.ContainerInfoRequest{NumStats: 1})
	if err == nil && ci != nil && len(ci.Stats) > 0 {
		return convertContainerInfoToMachineStats(ci.Stats[0])
	}

	return nil
}

func (ch *machineStatsFeed) Machine() *containers.MachineInfo {
	return ch.machine
}

func newMachineStatsFeed(manager manager.Manager, machine *containers.MachineInfo, root *containers.ContainerInfo) *machineStatsFeed {
	return &machineStatsFeed{
		manager: manager,
		root:    root,
		machine: machine,
	}
}

func newStatsChannel(manager manager.Manager, container *containers.ContainerInfo) *statsChannel {
	return &statsChannel{
		manager:   manager,
		container: container,
		closed:    false,
		ch:        make(chan *containers.Stats),
		doneCh:    make(chan bool),
	}
}

func calculateCpuUsage(nanoCpuTime uint64, numCores uint64) float64 {
	// https://github.com/kubernetes/heapster/issues/650
	return float64(nanoCpuTime) / float64(numCores*1e+9)
}

func convertContainerInfoToMachineStats(stats *v1.ContainerStats) *containers.MachineStats {
	cs := convertContainerInfoToStats(stats)
	return &containers.MachineStats{
		Cpu:    cs.Cpu,
		Memory: cs.Memory,
	}
}

func convertContainerInfoToStats(stats *v1.ContainerStats) *containers.Stats {
	cpu := &containers.CpuStats{
		TotalUsage:   stats.Cpu.Usage.Total,
		PerCoreUsage: stats.Cpu.Usage.PerCpu[:],
	}

	disk := &containers.DiskStats{
		PerDiskIo: make([]uint64, 0),
	}

	for _, d := range stats.DiskIo.IoServiceBytes {
		total := uint64(0)
		for _, bytesTransferred := range d.Stats {
			total += bytesTransferred
		}

		disk.PerDiskIo = append(disk.PerDiskIo, total)
	}

	net := &containers.NetworkStats{
		TotalRxBytes: stats.Network.RxBytes,
		TotalTxBytes: stats.Network.TxBytes,
		Interfaces:   make([]*containers.InterfaceStats, 0),
	}

	for _, nic := range stats.Network.Interfaces {
		net.Interfaces = append(net.Interfaces, &containers.InterfaceStats{
			Name:    nic.Name,
			RxBytes: nic.RxBytes,
			TxBytes: nic.TxBytes,
		})
	}

	return &containers.Stats{
		Cpu:     cpu,
		Memory:  &containers.MemoryStats{Usage: stats.Memory.Usage},
		Disk:    disk,
		Network: net,
	}
}
