package disk

import "zlnew/monitor-agent/internal/core"

type Collector struct {
	totalBytes float64
	freeBytes  float64
}

type DiskMetric = core.DiskMetric
