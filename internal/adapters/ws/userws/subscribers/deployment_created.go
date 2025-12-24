package subscribers

import (
	"horizonx-server/internal/adapters/ws/userws"
	"horizonx-server/internal/domain"
)

type DeploymentCreated struct {
	hub *userws.Hub
}

func NewDeploymentCreated(hub *userws.Hub) *DeploymentCreated {
	return &DeploymentCreated{hub: hub}
}

func (s *DeploymentCreated) Handle(event any) {
	evt, ok := event.(domain.EventDeploymentCreated)
	if !ok {
		return
	}

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: "deployments",
		Event:   "deployment_created",
		Payload: evt,
	})
}
