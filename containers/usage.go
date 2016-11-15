package containers

type MemoryUsage struct {
	// bytes used
	Bytes uint64 `json:"bytes"`
}

type InterfaceUsage struct {
	Name    string `json:"name"`
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}

type NetworkUsage struct {
	// total bytes received
	TotalRxBytes uint64 `json:"total_rx_bytes"`

	// total bytes sent
	TotalTxBytes uint64 `json:"total_tx_bytes"`

	// per interface stats
	Interfaces []*InterfaceUsage `json:"interfaces"`
}

type CpuUsage struct {
	// total CPU usage in nanoseconds
	Total int64 `json:"total"`

	// per core usage in nanoseconds
	PerCore []int64 `json:"per_core,omitempty"`
}

type DiskUsage struct {
	// per disk IO in bytes
	PerDiskIo []uint64 `json:"per_disk_io_bytes,omitempty"`
}

type Usage struct {
	Cpu     *CpuUsage     `json:"cpu,omitempty"`
	Memory  *MemoryUsage  `json:"memory,omitempty"`
	Network *NetworkUsage `json:"network,omitempty"`
	Disk    *DiskUsage    `json:"disk,omitempty"`
}

type MachineUsage struct {
	Cpu    *CpuUsage    `json:"cpu,omitempty"`
	Memory *MemoryUsage `json"memory,omitempty"`
}

type UsageChannel interface {
	Container() *ContainerInfo
	GetChannel() <-chan *Usage
	Close() error
}

type MachineUsageFeed interface {
	Machine() *MachineInfo
	Next() *MachineUsage
}
