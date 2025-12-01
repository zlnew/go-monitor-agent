// Package cpu
package cpu

import (
	"context"
)

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (any, error) {
	spec, _ := readSpec()
	usage, perCore, _ := readUsage()
	watt, _ := c.readWatt()
	temp, _ := readTemp()
	freq, _ := readFreq()

	return CPUMetric{
		Spec:      spec,
		Usage:     usage,
		PerCore:   perCore,
		Watt:      watt,
		Temp:      temp,
		Frequency: freq,
	}, nil
}
