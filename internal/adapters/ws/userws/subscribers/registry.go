package subscribers

import (
	"horizonx-server/internal/adapters/ws/userws"
)

func Register(bus EventBus, hub *userws.Hub) {
	serverStatusChanged := NewServerStatusChanged(hub)
	serverMetricsUpdated := NewServerMetricsUpdated(hub)

	bus.Subscribe("server_status_changed", serverStatusChanged.Handle)
	bus.Subscribe("server_metrics_updated", serverMetricsUpdated.Handle)
}
