// Package core
package core

type CPUMetric struct {
	Usage     float64   `json:"usage"`
	PerCore   []float64 `json:"per_core"`
	Watt      float64   `json:"watt"`
	Temp      float64   `json:"temp"`
	Frequency float64   `json:"frequency"`
}

type Metrics struct {
	CPU CPUMetric `json:"cpu"`
}
