package containers

type MemoryStats struct {
	// bytes used
	Usage uint64 `json:"usage"`
}

type InterfaceStats struct {
	Name    string `json:"name"`
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}

type NetworkStats struct {
	// total bytes received
	TotalRxBytes uint64 `json:"total_rx_bytes"`

	// total bytes sent
	TotalTxBytes uint64 `json:"total_tx_bytes"`

	// per interface stats
	Interfaces []*InterfaceStats `json:"interfaces"`
}

type CpuStats struct {
	// total CPU usage in nanoseconds
	TotalUsage uint64 `json:"total_usage"`

	// per core usage in nanoseconds
	PerCoreUsage []uint64 `json:"per_core_usage,omitempty"`
}

type DiskStats struct {
	// per disk IO in bytes
	PerDiskIo []uint64 `json:"per_disk_io_bytes,omitempty"`
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
