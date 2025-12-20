// Package application
package application

import (
	"context"
	"time"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/event"
	"horizonx-server/internal/logger"
)

type StatusListener struct {
	repo domain.ApplicationRepository
	log  logger.Logger
}

func NewStatusListener(repo domain.ApplicationRepository, log logger.Logger) *StatusListener {
	return &StatusListener{
		repo: repo,
		log:  log,
	}
}

func (l *StatusListener) Register(bus *event.Bus) {
	bus.Subscribe("application_status_changed", l.handleStatusChanged)
	bus.Subscribe("application_deployed", l.handleDeployed)
}

func (l *StatusListener) handleStatusChanged(event any) {
	evt, ok := event.(domain.EventApplicationStatusChanged)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := l.repo.UpdateStatus(ctx, evt.ApplicationID, evt.Status); err != nil {
		l.log.Error("failed to update application status",
			"app_id", evt.ApplicationID,
			"status", evt.Status,
			"error", err,
		)
	} else {
		l.log.Debug("application status updated",
			"app_id", evt.ApplicationID,
			"status", evt.Status,
		)
	}
}

func (l *StatusListener) handleDeployed(event any) {
	evt, ok := event.(domain.EventApplicationDeployed)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := l.repo.UpdateLastDeployment(ctx, evt.ApplicationID); err != nil {
		l.log.Error("failed to update application last deployment",
			"app_id", evt.ApplicationID,
			"error", err,
		)
	} else {
		l.log.Debug("application last deployment updated",
			"app_id", evt.ApplicationID,
		)
	}

	l.log.Info("application deployment completed",
		"app_id", evt.ApplicationID,
		"success", evt.Success,
	)
}
