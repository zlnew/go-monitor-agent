package memory

func (c *Collector) getMemTotal() float64 {
	return float64(c.memTotal) / 1024 / 1024 / 1024
}

func (c *Collector) getMemAvailable() float64 {
	return float64(c.memAvailable) / 1024 / 1024 / 1024
}

func (c *Collector) getMemUsed() float64 {
	used := c.memTotal - c.memAvailable
	return float64(used) / 1024 / 1024 / 1024
}
