package subscribers

import (
	"horizonx-server/internal/adapters/ws/userws"
	"horizonx-server/internal/domain"
)

type ServerStatusChanged struct {
	hub *userws.Hub
}

func NewServerStatusChanged(hub *userws.Hub) *ServerStatusChanged {
	return &ServerStatusChanged{hub: hub}
}

func (s *ServerStatusChanged) Handle(event any) {
	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: "server_status",
		Event:   "server_status_changed",
		Payload: event,
	})
}
