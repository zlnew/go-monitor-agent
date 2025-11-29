package core

import (
	"context"
	"time"

	"zlnew/monitor-agent/internal/infra/logger"
)

type Scheduler struct {
	reg      *Registry
	interval time.Duration
	log      logger.Logger
}

func NewScheduler(reg *Registry, interval time.Duration, log logger.Logger) *Scheduler {
	return &Scheduler{reg, interval, log}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.collect(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) collect(ctx context.Context) {
	s.reg.mu.RLock()
	collectors := s.reg.collectors
	s.reg.mu.RUnlock()

	for name, c := range collectors {
		val, err := c.Collect(ctx)
		if err == nil {
			s.reg.Update(name, val)
		} else {
			s.log.Error("collector", "name", name, "error", err)
		}
	}
}
