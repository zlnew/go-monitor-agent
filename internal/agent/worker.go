// Package agent
package agent

import (
	"context"
	"fmt"
	"time"

	"horizonx-server/internal/agent/executor"
	"horizonx-server/internal/config"
	"horizonx-server/internal/domain"
	"horizonx-server/internal/event"
	"horizonx-server/internal/logger"
)

type JobWorker struct {
	cfg      *config.Config
	log      logger.Logger
	executor *executor.Executor
	client   *Client
}

func NewJobWorker(cfg *config.Config, log logger.Logger, workDir string) *JobWorker {
	return &JobWorker{
		cfg:      cfg,
		log:      log,
		executor: executor.NewExecutor(log, workDir),
		client:   NewClient(cfg),
	}
}

func (w *JobWorker) Initialize() error {
	return w.executor.Initialize()
}

func (w *JobWorker) Start(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	w.log.Info("job worker started, polling for jobs...")

	for {
		select {
		case <-ctx.Done():
			w.log.Info("job worker stopping...")
			return ctx.Err()

		case <-ticker.C:
			if err := w.pollAndExecuteJobs(ctx); err != nil {
				w.log.Error("failed to poll jobs", "error", err)
			}
		}
	}
}

func (w *JobWorker) pollAndExecuteJobs(ctx context.Context) error {
	jobs, err := w.client.GetPendingJobs(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch jobs: %w", err)
	}

	if len(jobs) == 0 {
		return nil
	}

	w.log.Info("received jobs", "count", len(jobs))

	for _, job := range jobs {
		if err := w.processJob(ctx, job); err != nil {
			w.log.Error("failed to process job", "job_id", job.ID, "error", err)
		}
	}

	return nil
}

func (w *JobWorker) processJob(ctx context.Context, job domain.Job) error {
	w.log.Info("processing job", "id", job.ID, "type", job.Type)

	if err := w.client.StartJob(ctx, job.ID); err != nil {
		w.log.Error("failed to mark job as running", "job_id", job.ID, "error", err)
		return err
	}

	execErr := w.execute(ctx, &job)

	status := domain.JobSuccess
	if execErr != nil {
		status = domain.JobFailed
		w.log.Error("job execution failed", "job_id", job.ID, "error", execErr)
	} else {
		w.log.Info("job executed successfully", "job_id", job.ID)
	}

	if err := w.client.FinishJob(ctx, job.ID, status); err != nil {
		w.log.Error("failed to mark job as finished", "job_id", job.ID, "error", err)
		return err
	}

	return execErr
}

func (w *JobWorker) execute(ctx context.Context, job *domain.Job) error {
	bus := event.New()

	bus.Subscribe("log_emitted", func(e any) {
		evt := e.(domain.EventLogEmitted)
		if err := w.client.SendLog(ctx, &domain.LogEmitRequest{
			Timestamp:     evt.Timestamp,
			Level:         evt.Level,
			Source:        evt.Source,
			Action:        evt.Action,
			TraceID:       job.TraceID,
			JobID:         &job.ID,
			ServerID:      &job.ServerID,
			ApplicationID: job.ApplicationID,
			DeploymentID:  job.DeploymentID,
			Message:       evt.Message,
			Context:       evt.Context,
		}); err != nil {
			w.log.Error("failed to send log: %w", err)
		}
	})

	bus.Subscribe("commit_info_emitted", func(e any) {
		evt := e.(domain.EventCommitInfoEmitted)
		if err := w.client.SendCommitInfo(
			ctx,
			evt.DeploymentID,
			evt.Hash,
			evt.Message,
		); err != nil {
			w.log.Error("failed to send commit info: %w", err)
		}
	})

	onEmit := func(e any) {
		switch e.(type) {
		case domain.EventLogEmitted:
			bus.Publish("log_emitted", e)
		case domain.EventCommitInfoEmitted:
			bus.Publish("commit_info_emitted", e)
		}
	}

	return w.executor.Execute(ctx, job, onEmit)
}
