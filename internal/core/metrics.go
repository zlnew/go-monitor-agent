// Package core
package core

type CPUMetric struct {
	Spec      CPUSpec   `json:"spec"`
	Usage     float64   `json:"usage"`
	PerCore   []float64 `json:"per_core"`
	Watt      float64   `json:"watt"`
	Temp      float64   `json:"temp"`
	Frequency float64   `json:"frequency"`
}

type CPUSpec struct {
	Vendor    string `json:"vendor"`
	ModelName string `json:"model_name"`
	Cores     int    `json:"cores"`
	Threads   int    `json:"threads"`
}

type MemoryMetric struct {
	MemTotal     float64 `json:"mem_total"`
	MemAvailable float64 `json:"mem_available"`
	MemUsed      float64 `json:"mem_used"`
	SwapTotal    float64 `json:"swap_total"`
	SwapFree     float64 `json:"swap_free"`
	SwapUsed     float64 `json:"swap_used"`
}

type DiskMetric struct {
	Total float64 `json:"total"`
	Free  float64 `json:"free"`
	Used  float64 `json:"used"`
	Temp  float64 `json:"temp"`
}

type Metrics struct {
	CPU    CPUMetric    `json:"cpu"`
	Memory MemoryMetric `json:"memory"`
	Disk   DiskMetric   `json:"disk"`
}
