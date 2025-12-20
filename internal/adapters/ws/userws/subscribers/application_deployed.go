package subscribers

import (
	"fmt"

	"horizonx-server/internal/adapters/ws/userws"
	"horizonx-server/internal/domain"
)

type ApplicationDeployed struct {
	hub *userws.Hub
}

func NewApplicationDeployed(hub *userws.Hub) *ApplicationDeployed {
	return &ApplicationDeployed{hub: hub}
}

func (s *ApplicationDeployed) Handle(event any) {
	evt, ok := event.(domain.EventApplicationDeployed)
	if !ok {
		return
	}

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: fmt.Sprintf("application:%d", evt.ApplicationID),
		Event:   "application_deployed",
		Payload: evt,
	})

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: "applications",
		Event:   "application_deployed",
		Payload: evt,
	})
}
