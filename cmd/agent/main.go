package main

import (
	"context"

	"zlnew/monitor-agent/internal/agent"
	"zlnew/monitor-agent/internal/infra/config"
	"zlnew/monitor-agent/internal/infra/logger"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()
	log := logger.New(cfg)

	a := agent.New(log, cfg)
	if err := a.Run(ctx); err != nil {
		log.Fatal("agent stopped with error:", err)
	}
}
