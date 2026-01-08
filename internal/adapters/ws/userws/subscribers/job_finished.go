package subscribers

import (
	"fmt"

	"horizonx/internal/adapters/ws/userws"
	"horizonx/internal/domain"
)

type JobFinished struct {
	hub *userws.Hub
}

func NewJobFinished(hub *userws.Hub) *JobFinished {
	return &JobFinished{hub: hub}
}

func (s *JobFinished) Handle(event any) {
	evt, ok := event.(domain.EventJobFinished)
	if !ok {
		return
	}

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: fmt.Sprintf("job:%d", evt.JobID),
		Event:   "job_finished",
		Payload: evt,
	})

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: "jobs",
		Event:   "job_finished",
		Payload: evt,
	})
}
