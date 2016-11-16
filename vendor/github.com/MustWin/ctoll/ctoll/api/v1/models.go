package v1

import (
	"encoding/json"
	"net/http"
)

type CreateOrgRequest struct {
	Name          string `json:"name"`
	BillingPlanID string `json:"billing_plan_id"`
}

type Org struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type APIKey struct {
	Key   string `json:"key"`
	OrgID string `json:"org_id"`
}

type MachineInfo struct {
	SystemUuid      string            `json:"system_uuid"`
	Cores           int               `json:"cores"`
	MemoryBytes     uint64            `json:"memory_bytes"`
	CpuFrequencyKhz uint64            `json:"cpu_frequency_khz"`
	Labels          map[string]string `json:"labels"`
	Name            string            `json:"name"`
}

type ClusterCapacity struct {
	CPUShares   float64 `json:"cpu_shares"`
	MemoryBytes int64   `json:"memory_bytes"`
	NumNodes    int     `json:"node_count"`
}

type MachineUsage struct {
	CPUShares   float64 `json:"cpu_shares"`
	MemoryBytes int64   `json:"memory_bytes"`
}

type Usage struct {
	// nanoseconds of CPU time
	CPU int64 `json:"cpu_time"`

	// shares of CPU (not used by cmeter)
	CPUShares float64 `json:"cpu_shares,omitempty"`

	// memory in bytes
	MemoryBytes int64 `json:"memory_bytes"`

	// disk in bytes
	DiskIOBytes int64 `json:"disk_io_bytes"`

	// network recv in bytes
	NetworkRxBytes int64 `json:"net_rx_bytes"`

	// network sent in bytes
	NetworkTxBytes int64 `json:"net_tx_bytes"`
}

func (u *Usage) Add(u2 *Usage) *Usage {
	return &Usage{
		CPU:            u.CPU + u2.CPU,
		CPUShares:      u.CPUShares + u2.CPUShares,
		MemoryBytes:    u.MemoryBytes + u2.MemoryBytes,
		DiskIOBytes:    u.DiskIOBytes + u2.DiskIOBytes,
		NetworkRxBytes: u.NetworkRxBytes + u2.NetworkRxBytes,
		NetworkTxBytes: u.NetworkTxBytes + u2.NetworkTxBytes,
	}
}

func (u *Usage) Average(n int64) *Usage {
	return &Usage{
		CPU:            u.CPU / n,
		CPUShares:      u.CPUShares / float64(n),
		MemoryBytes:    u.MemoryBytes / n,
		DiskIOBytes:    u.DiskIOBytes / n,
		NetworkRxBytes: u.NetworkRxBytes / n,
		NetworkTxBytes: u.NetworkTxBytes / n,
	}
}

type BlockAlloc struct {
	CPUShares   float64 `json:"cpu_shares"`
	MemoryBytes int64   `json:"memory_bytes"`
}

type ContainerInfo struct {
	ImageName string            `json:"image_name,omitempty"`
	ImageTag  string            `json:"image_tag,omitempty"`
	Name      string            `json:"name,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	Machine   *MachineInfo      `json:"machine,omitempty"`
}

type MeterEventType string

var (
	MeterEventTypeStart         MeterEventType = "start"
	MeterEventTypeStop          MeterEventType = "stop"
	MeterEventTypeSample        MeterEventType = "sample"
	MeterEventTypeMachineSample MeterEventType = "machine_sample"
)

type MeterEvent struct {
	MeterID   string         `json:"meter_id"`
	Type      MeterEventType `json:"event_type"`
	Timestamp int64          `json:"timestamp"`
}

type OrgCreateRequest struct {
	Name string `json:"name"`
}

type StartMeterEvent struct {
	*MeterEvent
	Container *ContainerInfo `json:"container"`
	Allocated *BlockAlloc    `json:"allocated"`
}

type StopMeterEvent struct {
	*MeterEvent
	Container *ContainerInfo `json:"container"`
}

type SampleMeterEvent struct {
	*MeterEvent
	Container *ContainerInfo `json:"container"`
	Usage     *Usage         `json:"usage"`
}

type MachineSampleMeterEvent struct {
	*MeterEvent
	Machine *MachineInfo  `json:"machine"`
	Usage   *MachineUsage `json:"usage"`
}

type SessionState string

var (
	StateNone   SessionState
	StateActive SessionState = "active"
	StateClosed SessionState = "closed"
)

type MeterSample struct {
	SessionID string `json:"session_id"`
	OrgID     string `json:"org_id"`
	Timestamp int64  `json:"timestamp"`
	Usage     *Usage `json:"usage"`
}

type MeterSession struct {
	ID        string       `json:"id"`
	OrgID     string       `json:"org_id"`
	APIKey    string       `json:"apikey"`
	MeterID   string       `json:"meter_id"`
	StartTime int64        `json:"start_time"`
	EndTime   int64        `json:"end_time"`
	State     SessionState `json:"state"`
	Machine   *MachineInfo `json:"machine"`

	SampleCount  int32          `json:"sample_count"`
	Container    *ContainerInfo `json:"container"`
	Allocated    *BlockAlloc    `json:"allocated"`
	TotalUsage   *Usage         `json:"total_usage"`
	AverageUsage *Usage         `json:"average_usage"`
}

type SessionStateTransition struct {
	SessionID string       `json:"session_id"`
	OrgID     string       `json:"org_id"`
	MeterID   string       `json:"meter_id"`
	State     SessionState `json:"state"`
	Timestamp int64        `json:"timestamp"`
	Allocated *BlockAlloc  `json:"allocated"`
	Machine   *MachineInfo `json:"machine"`
}

type SessionRecommendation struct {
	*MeterSession
	Recommendation string `json:"recommendation"`
}

func ServeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
