# go-monitor-agent

Lightweight Linux metrics agent that scrapes `/proc` and `/sys` once per second and exposes a JSON snapshot at `/metrics`.

## Requirements

- Go 1.25.4+.
- Linux host with `/proc` and `/sys` available.
- Optional kernel interfaces/sensors:
  - RAPL or hwmon power/temperature files for CPU and GPU.
  - `/sys/class/drm/*/gpu_busy_percent` and `mem_info_vram_*` for GPU utilization and VRAM.
  - NVMe hwmon temperature for disk temperature.
- Missing files are tolerated; collectors will emit zero values instead of failing.

## Run, build, test

- Run locally: `go run ./cmd/agent` (binds to `:3000`).
- Configure address: set `HTTP_ADDR` in the env or a `.env` file, e.g. `HTTP_ADDR=:8080`.
- Build: `go build -o bin/monitor-agent ./cmd/agent` or `make build`.
- Clean build artifact: `make clean`.
- Tests: `go test ./...` (none yet).

## Metrics

Collectors refresh every second and degrade gracefully when inputs are unreadable.

- CPU: usage and per-core usage percent (averaged since boot from `/proc/stat`), temperature via hwmon CPU sensors, current frequency (MHz) from cpufreq, and `power_watt` from RAPL or hwmon when available.
- GPU: first detected card; usage from `/sys/class/drm/*/gpu_busy_percent`, VRAM total/used in GB from `mem_info_vram_*`, temperature/power from hwmon GPU nodes; returned as a slice even if there is only one GPU.
- Memory: total, available, used, and swap total/free/used in GiB from `/proc/meminfo`.
- Disk: root filesystem total/free/used in decimal GB from `statfs` and NVMe temperature from hwmon if present.
- Network: aggregate of non-loopback interfaces from `/proc/net/dev` with cumulative `rx_bytes`/`tx_bytes` and Mbps rates computed between samples (first sample reports `0` speeds).
- Uptime: `uptime_seconds` from `/proc/uptime`.

### Endpoint

`GET /metrics` returns the latest snapshot:

```json
{
  "cpu": {
    "usage": 12.5,
    "per_core": [8.2, 14.1, 13.3, 9.9],
    "temperature": 55.0,
    "frequency": 2200.0,
    "power_watt": 7.8
  },
  "gpu": [
    {
      "id": 0,
      "vendor": "AMD",
      "model": "RX 6600",
      "temperature": 48.3,
      "usage": 15.7,
      "vram_total_gb": 8.0,
      "vram_used_gb": 2.1,
      "vram_usage_percent": 0,
      "power_watt": 65.2,
      "fan_speed_percent": 0
    }
  ],
  "memory": {
    "total_gb": 15.8,
    "used_gb": 6.4,
    "available_gb": 9.4,
    "swap_total_gb": 2.0,
    "swap_free_gb": 2.0,
    "swap_used_gb": 0.0
  },
  "disk": [
    {
      "name": "nvme",
      "total_gb": 476.3,
      "used_gb": 174.2,
      "free_gb": 302.1,
      "percent": 56,
      "temperature": 41.5
    }
  ],
  "network": {
    "rx_bytes": 123456789,
    "tx_bytes": 987654321,
    "rx_speed": 1.2,
    "tx_speed": 8.5
  },
  "uptime_seconds": 12345
}
```

## Notes

- The HTTP server runs until process exit; the scheduler stops when context is canceled.
- Build artifacts land in `bin/`; avoid committing binaries unless intentional.
