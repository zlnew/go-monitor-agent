// Package http
package http

import (
	"context"
	"encoding/json"
	"net/http"

	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/config"
	"zlnew/monitor-agent/internal/infra/logger"
)

type Server struct {
	reg *core.Registry
	log logger.Logger
	cfg *config.Config
}

func NewServer(cfg *config.Config, reg *core.Registry, log logger.Logger) *Server {
	return &Server{cfg: cfg, reg: reg, log: log}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", s.handleMetrics)

	s.log.Info("starting http server on " + s.cfg.Address)
	return http.ListenAndServe(s.cfg.Address, mux)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	data := s.reg.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
