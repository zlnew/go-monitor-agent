// Package cpu
package cpu

import (
	"context"
)

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (CPUMetric, error) {
	usage, perCore, _ := readUsage()
	watt, _ := readWatt()
	temp, _ := readTemp()
	freq, _ := readFreq()

	return CPUMetric{
		Usage:     usage,
		PerCore:   perCore,
		Watt:      watt,
		Temp:      temp,
		Frequency: freq,
	}, nil
}
