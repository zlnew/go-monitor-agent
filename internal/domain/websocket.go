package domain

import "fmt"

const (
	ChannelServerStatus          = "server_status"
	ChannelServerMetricsTemplate = "server:%d:metrics"
)

const (
	EventServerStatusUpdated   = "server_status_updated"
	EventServerMetricsReport   = "server_metrics_report"
	EventServerMetricsReceived = "server_metrics_received"
)

type ServerStatusPayload struct {
	ServerID int64 `json:"server_id"`
	IsOnline bool  `json:"is_online"`
}

func GetServerMetricsChannel(serverID int64) string {
	return fmt.Sprintf(ChannelServerMetricsTemplate, serverID)
}
