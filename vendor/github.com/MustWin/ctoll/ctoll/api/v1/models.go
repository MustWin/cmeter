package v1

import (
	"encoding/json"
	"net/http"
)

type MeterBilling struct {
	Period int64   `json:"period"`
	Price  float64 `json:"price"`
	Unit   float64 `json:"unit"`
}

type AllocationBilling struct {
	Price float64 `json:"price"`
	Unit  float64 `json:"unit"`
}

type BillingModel struct {
	OrgID       string             `json:"org_id"`
	CPU         *MeterBilling      `json:"cpu"`
	CPUAlloc    *AllocationBilling `json:"cpu_alloc"`
	Memory      *MeterBilling      `json:"memory"`
	MemoryAlloc *AllocationBilling `json:"memory_alloc"`
	DiskIO      *AllocationBilling `json:"disk_io"`
	NetworkRx   *AllocationBilling `json:"net_rx"`
	NetworkTx   *AllocationBilling `json:"net_tx"`
}

type Org struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type APIKey struct {
	Key   string `json:"key"`
	OrgID string `json:"org_id"`
}

type Usage struct {
	TotalCPUPerc   float32 `json:"total_cpu_perc"`
	MemoryBytes    int64   `json:"memory_bytes"`
	DiskIOBytes    int64   `json:"disk_io_bytes"`
	NetworkRxBytes int64   `json:"net_rx_bytes"`
	NetworkTxBytes int64   `json:"net_tx_bytes"`
}

type BlockAlloc struct {
	MaxCPUPerc  float32 `json:"max_cpu_perc"`
	MemoryBytes int64   `json:"memory_bytes"`
}

type ContainerInfo struct {
	ImageName string            `json:"image_name"`
	ImageTag  string            `json:"image_tag"`
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
}

type MeterEventType string

var (
	MeterEventTypeStart  MeterEventType = "start"
	MeterEventTypeStop   MeterEventType = "stop"
	MeterEventTypeSample MeterEventType = "sample"
)

type MeterEvent struct {
	MeterID   string         `json:"meter_id"`
	Type      MeterEventType `json:"event_type"`
	Timestamp int64          `json:"timestamp"`
	Container *ContainerInfo `json:"container"`
}

type StartMeterEvent struct {
	*MeterEvent
	Allocated *BlockAlloc `json:"allocated"`
}

type StopMeterEvent struct {
	*MeterEvent
}

type SampleMeterEvent struct {
	*MeterEvent
	Usage *Usage `json:"usage"`
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

	SampleCount  int32          `json:"sample_count"`
	Container    *ContainerInfo `json:"container"`
	Allocated    *BlockAlloc    `json:"allocated"`
	TotalUsage   *Usage         `json:"total_usage"`
	AverageUsage *Usage         `json:"average_usage"`
}

func ServeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
