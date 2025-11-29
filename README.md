# go-monitor-agent

A lightweight Go agent that exposes local CPU metrics over HTTP as JSON. It polls `/proc` and `/sys` to collect usage, per-core load, wattage, temperature, and frequency, then serves a simple `/metrics` endpoint that other systems can scrape.

## Getting started

- Prerequisite: Go 1.25 (per `go.mod`).
- Run locally: `go run ./cmd/agent` (binds to `:3000` by default).
- Configure address: set `HTTP_ADDR`, either in the environment or a `.env` file (e.g., `HTTP_ADDR=:8080`).
- Build a binary: `go build -o bin/agent ./cmd/agent`.
- Run tests: `go test ./...` (no tests yet).

## HTTP endpoint

- `GET /metrics` returns the latest snapshot.
- Example response:
  ```json
  {
    "cpu": {
      "usage": 0.42,
      "per_core": [0.35, 0.51, 0.41, 0.42],
      "watt": 5.8,
      "temp": 49.2,
      "frequency": 2.4
    }
  }
  ```

## Project layout

- `cmd/agent`: Entrypoint that wires config/logger and launches the agent.
- `internal/agent`: Builds the registry, registers collectors, and starts the HTTP server.
- `internal/core`: Metric types, registry, and scheduler that refreshes snapshots.
- `internal/collector/cpu`: CPU collector reading `/proc` and `/sys` for usage, wattage, temperature, and frequency.
- `internal/transport/http`: HTTP server exposing `/metrics` via the configured address.
- `internal/infra`: Config loader (`HTTP_ADDR`, `.env`) and logger interface/implementation.
- `pkg`: Small shared helpers.
- `bin/`: Built artifacts; avoid committing new binaries unless intentional.

## Notes

- Metrics are refreshed on a 1s tick by the scheduler; collectors attempt to degrade gracefully if expected files are missing.
- Logs use the standard library logger via the provided interface.
