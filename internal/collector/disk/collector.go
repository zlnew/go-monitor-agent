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

	total := c.getTotal()
	free := c.getFree()
	used := c.getUsed()
	temp, err := getTemp()
	if err != nil {
		return DiskMetric{}, err
	}

	return DiskMetric{
		Total: total,
		Free:  free,
		Used:  used,
		Temp:  temp,
	}, nil
}
