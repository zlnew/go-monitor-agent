package subscribers

import (
	"fmt"

	"horizonx/internal/adapters/ws/userws"
	"horizonx/internal/domain"
)

type ApplicationStatusChanged struct {
	hub *userws.Hub
}

func NewApplicationStatusChanged(hub *userws.Hub) *ApplicationStatusChanged {
	return &ApplicationStatusChanged{hub: hub}
}

func (s *ApplicationStatusChanged) Handle(event any) {
	evt, ok := event.(domain.EventApplicationStatusChanged)
	if !ok {
		return
	}

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: fmt.Sprintf("application:%d", evt.ApplicationID),
		Event:   "application_status_changed",
		Payload: evt,
	})

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: "applications",
		Event:   "application_status_changed",
		Payload: evt,
	})
}
