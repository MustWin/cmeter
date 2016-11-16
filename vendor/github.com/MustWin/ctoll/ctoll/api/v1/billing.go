package v1

import ()

// TODO: move all this billing calculation stuff into its own "billing" package but keep models here
type MeterBilling struct {
	Period int64   `json:"period"`
	Price  float64 `json:"price"`
	Unit   float64 `json:"unit"`
}

func (b *MeterBilling) Enabled() bool {
	return b.Period > 0 && b.Price > 0 && b.Unit > 0
}

type AllocationBilling struct {
	Price float64 `json:"price"`
	Unit  float64 `json:"unit"`
}

func (b *AllocationBilling) Enabled() bool {
	return b.Price > 0 && b.Unit > 0
}

type BillingModel struct {
	OrigPlanID  string             `json:"orig_plan_id"`
	PlanID      string             `json:"plan_id"`
	OrgID       string             `json:"org_id"`
	BasePrice   float64            `json:"base_price"`
	CPU         *MeterBilling      `json:"cpu"`
	CPUAlloc    *AllocationBilling `json:"cpu_alloc"`
	Memory      *MeterBilling      `json:"memory"`
	MemoryAlloc *AllocationBilling `json:"memory_alloc"`
	DiskIO      *AllocationBilling `json:"disk_io"`
	NetworkRx   *AllocationBilling `json:"net_rx"`
	NetworkTx   *AllocationBilling `json:"net_tx"`
}

type BillingPlan struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	BasePrice   float64            `json:"base_price"`
	CPU         *MeterBilling      `json:"cpu"`
	CPUAlloc    *AllocationBilling `json:"cpu_alloc"`
	Memory      *MeterBilling      `json:"memory"`
	MemoryAlloc *AllocationBilling `json:"memory_alloc"`
	DiskIO      *AllocationBilling `json:"disk_io"`
	NetworkRx   *AllocationBilling `json:"net_rx"`
	NetworkTx   *AllocationBilling `json:"net_tx"`
}

func (p *BillingPlan) Model() *BillingModel {
	return &BillingModel{
		OrigPlanID:  p.ID,
		BasePrice:   p.BasePrice,
		CPU:         p.CPU,
		CPUAlloc:    p.CPUAlloc,
		Memory:      p.Memory,
		MemoryAlloc: p.MemoryAlloc,
		DiskIO:      p.DiskIO,
		NetworkRx:   p.NetworkRx,
		NetworkTx:   p.NetworkTx,
	}
}

// NOTE: don't love this name, alternatives welcome
type DistributionInfo struct {
	OrgID   string                `json:"org_id"`
	OrgName string                `json:"org_name"`
	Data    []*PeriodDistribution `json:"data"`
}

// NOTE: don't love this name, alternatives welcome
type PeriodDistribution struct {
	Start              int64   `json:"start"`
	End                int64   `json:"end"`
	AverageCPU         float64 `json:"average_cpu"`
	AverageCPUAlloc    float64 `json:"average_cpu_alloc"`
	AverageMemory      int64   `json:"average_memory"`
	AverageMemoryAlloc int64   `json:"average_memory_alloc"`
	DiskIO             int64   `json:"disk_io"`
	NetRx              int64   `json:"net_rx"`
	NetTx              int64   `json:"net_tx"`
	Cost               float64 `json:"cost"`
}

func (d *PeriodDistribution) AddUsage(u *Usage) {
	d.AverageCPU += u.CPUShares
	d.AverageMemory += u.MemoryBytes
	d.DiskIO += u.DiskIOBytes
	d.NetRx += u.NetworkRxBytes
	d.NetTx += u.NetworkTxBytes
}

func (d *PeriodDistribution) AddAlloc(b *BlockAlloc) {
	d.AverageCPUAlloc += b.CPUShares
	d.AverageMemoryAlloc += b.MemoryBytes
}

func (d *PeriodDistribution) AverageUsageOver(n int64) {
	if n == 0 {
		return
	}

	f := float64(n)
	d.AverageCPU = d.AverageCPU / f
	d.AverageMemory = d.AverageMemory / int64(n)
}

func (d *PeriodDistribution) AverageAllocationOver(n int64) {
	if n == 0 {
		return
	}

	f := float64(n)
	d.AverageCPUAlloc = d.AverageCPUAlloc / f
	d.AverageMemoryAlloc = d.AverageMemoryAlloc / n
}
