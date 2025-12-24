package subscribers

import (
	"horizonx-server/internal/adapters/ws/userws"
	"horizonx-server/internal/domain"
)

type ApplicationCreated struct {
	hub *userws.Hub
}

func NewApplicationCreated(hub *userws.Hub) *ApplicationCreated {
	return &ApplicationCreated{hub: hub}
}

func (s *ApplicationCreated) Handle(event any) {
	evt, ok := event.(domain.EventApplicationCreated)
	if !ok {
		return
	}

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: "applications",
		Event:   "application_created",
		Payload: evt,
	})
}
