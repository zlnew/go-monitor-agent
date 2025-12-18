package subscribers

import (
	"horizonx-server/internal/adapters/ws/userws"
	"horizonx-server/internal/domain"
)

type ServerMetricsUpdated struct {
	hub *userws.Hub
}

func NewServerMetricsUpdated(hub *userws.Hub) *ServerMetricsUpdated {
	return &ServerMetricsUpdated{hub: hub}
}

func (s *ServerMetricsUpdated) Handle(event any) {
	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: "server_metrics",
		Event:   "server_metrics_updated",
		Payload: event,
	})
}
