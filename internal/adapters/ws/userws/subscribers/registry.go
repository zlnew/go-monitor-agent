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

	// Job Events
	jobCreated := NewJobCreated(hub)
	jobStatusChanged := NewJobStatusChanged(hub)

	bus.Subscribe("job_created", jobCreated.Handle)
	bus.Subscribe("job_status_changed", jobStatusChanged.Handle)

	// Application Events
	applicationCreated := NewApplicationCreated(hub)
	applicationStatusChanged := NewApplicationStatusChanged(hub)

	bus.Subscribe("application_created", applicationCreated.Handle)
	bus.Subscribe("application_status_changed", applicationStatusChanged.Handle)

	// Deployment Events
	deploymentCreated := NewDeploymentCreated(hub)
	deploymentStatusChanged := NewDeploymentStatusChanged(hub)
	deploymentLogsUpdated := NewDeploymentLogsUpdated(hub)
	deploymentCommitInfoReceived := NewDeploymentCommitInfoReceived(hub)

	bus.Subscribe("deployment_created", deploymentCreated.Handle)
	bus.Subscribe("deployment_status_changed", deploymentStatusChanged.Handle)
	bus.Subscribe("deployment_logs_updated", deploymentLogsUpdated.Handle)
	bus.Subscribe("deployment_commit_info_received", deploymentCommitInfoReceived.Handle)
}
