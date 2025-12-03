// Package agent
package agent

import (
	"zlnew/monitor-agent/internal/collector/cpu"
	"zlnew/monitor-agent/internal/collector/disk"
	"zlnew/monitor-agent/internal/collector/gpu"
	"zlnew/monitor-agent/internal/collector/memory"
	"zlnew/monitor-agent/internal/collector/network"
	"zlnew/monitor-agent/internal/collector/uptime"
	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/config"
	"zlnew/monitor-agent/internal/infra/logger"
	"zlnew/monitor-agent/internal/transport/http"
)

type Agent struct {
	log  logger.Logger
	cfg  *config.Config
	reg  *core.Registry
	http *http.Server
}

func New(log logger.Logger, cfg *config.Config) *Agent {
	reg := core.NewRegistry()
	cpuCollector := cpu.NewCollector()
	gpuCollector := gpu.NewCollector()
	memoryCollector := memory.NewCollector()
	diskCollector := disk.NewCollector()
	networkCollector := network.NewCollector()
	uptimeCollector := uptime.NewCollector()

	reg.Register("cpu", cpuCollector)
	reg.Register("gpu", gpuCollector)
	reg.Register("memory", memoryCollector)
	reg.Register("disk", diskCollector)
	reg.Register("network", networkCollector)
	reg.Register("uptime", uptimeCollector)

	httpServer := http.NewServer(cfg, reg, log)

	return &Agent{
		log:  log,
		cfg:  cfg,
		reg:  reg,
		http: httpServer,
	}
}
