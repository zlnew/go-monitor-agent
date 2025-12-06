// Package agent
package agent

import (
	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/config"
	"zlnew/monitor-agent/internal/infra/logger"
	"zlnew/monitor-agent/internal/transport/http"
)

type Agent struct {
	log     logger.Logger
	cfg     *config.Config
	sampler *Sampler
	store   *core.SnapshotStore
	http    *http.Server
	hub     *http.Hub
}

func New(log logger.Logger, cfg *config.Config, hub *http.Hub) *Agent {
	sam := NewSampler(log)
	store := core.NewSnapshotStore()
	httpServer := http.NewServer(cfg, store, log, hub)

	return &Agent{
		log:     log,
		cfg:     cfg,
		sampler: sam,
		store:   store,
		http:    httpServer,
		hub:     hub,
	}
}
