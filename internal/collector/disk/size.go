package disk

func (c *Collector) getTotalGB() float64 {
	return c.totalBytes / 1_000_000_000
}

func (c *Collector) getFreeGB() float64 {
	return c.freeBytes / 1_000_000_000
}

func (c *Collector) getUsedGB() float64 {
	usedBytes := c.totalBytes - c.freeBytes
	return usedBytes / 1_000_000_000
}
