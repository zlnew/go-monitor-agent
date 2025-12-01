package memory

func (c *Collector) getSwapTotal() float64 {
	return float64(c.swapTotal) / 1024 / 1024 / 1024
}

func (c *Collector) getSwapFree() float64 {
	return float64(c.swapFree) / 1024 / 1024 / 1024
}

func (c *Collector) getSwapUsed() float64 {
	used := c.swapTotal - c.swapFree
	return float64(used) / 1024 / 1024 / 1024
}
