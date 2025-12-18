package subscribers

import (
	"horizonx-server/internal/adapters/ws/userws"
	"horizonx-server/internal/domain"
)

type ServerMetricsReceived struct {
	hub *userws.Hub
}

func NewServerMetricsReceived(hub *userws.Hub) *ServerMetricsReceived {
	return &ServerMetricsReceived{hub: hub}
}

func (s *ServerMetricsReceived) Handle(event any) {
	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: "server_metrics",
		Event:   "server_metrics_received",
		Payload: event,
	})
}
