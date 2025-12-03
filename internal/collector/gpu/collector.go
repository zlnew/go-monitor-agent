// Package gpu
package gpu

import "context"

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (any, error) {
	usage := getUsage()
	temperature := getTemperature()
	vramTotalGB := getVramTotalGB()
	vramUsedGB := getVramUsedGB()
	powerWatt := getPowerWatt()

	return []GPUMetric{{
		Usage:       usage,
		Temperature: temperature,
		VramTotalGB: vramTotalGB,
		VramUsedGB:  vramUsedGB,
		PowerWatt:   powerWatt,
	}}, nil
}
