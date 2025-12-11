package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"horizonx-server/internal/agent"
	"horizonx-server/internal/config"
	"horizonx-server/internal/core/metrics"
	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: No .env file found, relying on system environment variables")
	}

	serverURL := os.Getenv("HORIZONX_SERVER_URL")
	if serverURL == "" {
		serverURL = "ws://localhost:3000/ws"
	}

	agentToken := os.Getenv("HORIZONX_AGENT_TOKEN")
	if agentToken == "" {
		log.Fatal("FATAL: HORIZONX_AGENT_TOKEN is missing in .env or system vars!")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	appLog := logger.New(cfg)

	appLog.Info("HorizonX Agent: Starting spy mission...", "server_url", serverURL)

	agentClient := agent.NewAgent(serverURL, agentToken, appLog)

	metricsSampler := metrics.NewSampler(appLog)
	metricsSink := func(m domain.Metrics) { agentClient.SendMetric(m) }
	metricsScheduler := metrics.NewScheduler(cfg.Interval, appLog, metricsSampler.Collect, metricsSink)
	go metricsScheduler.Start(ctx)

	if err := agentClient.Run(ctx); err != nil && err != context.Canceled {
		appLog.Error("agent run failed unexpectedly", "error", err)
	}

	appLog.Info("Agent stopped gracefully.")
}
