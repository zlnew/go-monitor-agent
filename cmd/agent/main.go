package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

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
		serverURL = "http://localhost:3000/api/metrics/report"
	}

	agentToken := os.Getenv("HORIZONX_AGENT_TOKEN")
	if agentToken == "" {
		log.Fatal("FATAL: HORIZONX_AGENT_TOKEN is missing in .env or system vars!")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("HorizonX Agent: Starting spy mission...")
	log.Printf("Target Server: %s", serverURL)

	logCfg := &config.Config{
		LogLevel:  os.Getenv("LOG_LEVEL"),
		LogFormat: os.Getenv("LOG_FORMAT"),
	}

	appLog := logger.New(logCfg)
	sampler := metrics.NewSampler(appLog)
	client := &http.Client{Timeout: 5 * time.Second}

	dataSink := func(m domain.Metrics) {
		if err := sendMetrics(client, serverURL, agentToken, m); err != nil {
			appLog.Error("agent", "delivery_failed", err.Error())
		} else {
			appLog.Debug("agent", "status", "delivered")
		}
	}

	scheduler := metrics.NewScheduler(2*time.Second, appLog, sampler.Collect, dataSink)
	scheduler.Start(ctx)

	log.Println("Agent stopped gracefully.")
}

func sendMetrics(client *http.Client, url, token string, data domain.Metrics) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("create req: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("network: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("server rejected: %d", resp.StatusCode)
	}

	return nil
}
