package cpu

import "zlnew/monitor-agent/internal/core"

type Collector struct {
	lastEnergy uint64
	lastTime   int64
}

type CPUMetric = core.CPUMetric
