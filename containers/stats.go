package containers

type MemoryStats struct {
	UsageBytes uint64 `json:"usage"`
}

type InterfaceStats struct {
	Name    string `json:"name"`
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}

type NetworkStats struct {
	TotalRxBytes uint64            `json:"total_rx_bytes"`
	TotalTxBytes uint64            `json:"total_tx_bytes"`
	Interfaces   []*InterfaceStats `json:"interfaces"`
}

type CpuStats struct {
	TotalUsagePerc   float64   `json:"total_usage"`
	PerCoreUsagePerc []float64 `json:"per_core_usage,omitempty"`
}

type DiskStats struct {
	PerDiskIoBytes []uint64 `json:"per_disk_io_bytes,omitempty"`
}

type Stats struct {
	Cpu     *CpuStats     `json:"cpu,omitempty"`
	Memory  *MemoryStats  `json:"memory,omitempty"`
	Network *NetworkStats `json:"network,omitempty"`
	Disk    *DiskStats    `json:"disk,omitempty"`
}

type MachineStats struct {
	Cpu    *CpuStats    `json:"cpu,omitempty"`
	Memory *MemoryStats `json"memory,omitempty"`
}

type StatsChannel interface {
	Container() *ContainerInfo
	GetChannel() <-chan *Stats
	Close() error
}

type MachineStatsFeed interface {
	Machine() *MachineInfo
	Next() *MachineStats
}
