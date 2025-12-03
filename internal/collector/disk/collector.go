// Package disk
package disk

import "context"

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (any, error) {
	err := c.readStat()
	if err != nil {
		return DiskMetric{}, err
	}

	totalGB := c.getTotalGB()
	freeGB := c.getFreeGB()
	usedGB := c.getUsedGB()
	temperature := getTemperature()

	return []DiskMetric{{
		TotalGB:     totalGB,
		FreeGB:      freeGB,
		UsedGB:      usedGB,
		Temperature: temperature,
	}}, nil
}
