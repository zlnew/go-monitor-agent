package subscribers

import (
	"fmt"

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
	evt, ok := event.(domain.Metrics)
	if !ok {
		return
	}

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: fmt.Sprintf("server_metrics:%s", evt.ServerID.String()),
		Event:   "server_metrics_received",
		Payload: evt,
	})
}
