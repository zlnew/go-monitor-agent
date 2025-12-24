package subscribers

import (
	"horizonx-server/internal/adapters/ws/userws"
	"horizonx-server/internal/domain"
)

type JobCreated struct {
	hub *userws.Hub
}

func NewJobCreated(hub *userws.Hub) *JobCreated {
	return &JobCreated{hub: hub}
}

func (s *JobCreated) Handle(event any) {
	evt, ok := event.(domain.EventJobCreated)
	if !ok {
		return
	}

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: "jobs",
		Event:   "job_created",
		Payload: evt,
	})
}
