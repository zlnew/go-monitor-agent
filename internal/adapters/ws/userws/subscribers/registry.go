package subscribers

import (
	"horizonx-server/internal/adapters/ws/userws"
)

func Register(bus EventBus, hub *userws.Hub) {
	// Server Events
	serverStatusChanged := NewServerStatusChanged(hub)
	serverMetricsReceived := NewServerMetricsReceived(hub)

	bus.Subscribe("server_status_changed", serverStatusChanged.Handle)
	bus.Subscribe("server_metrics_received", serverMetricsReceived.Handle)

	// Application Events
	applicationStatusChanged := NewApplicationStatusChanged(hub)
	applicationDeployed := NewApplicationDeployed(hub)

	bus.Subscribe("application_status_changed", applicationStatusChanged.Handle)
	bus.Subscribe("application_deployed", applicationDeployed.Handle)

	// Job Events
	jobStatusChanged := NewJobStatusChanged(hub)

	bus.Subscribe("job_status_changed", jobStatusChanged.Handle)
}
