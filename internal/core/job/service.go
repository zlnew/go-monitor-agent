// Package job
package job

import (
	"context"

	"horizonx-server/internal/domain"
)

type JobService struct {
	repo domain.JobRepository
}

func NewService(repo domain.JobRepository) domain.JobService {
	return &JobService{repo: repo}
}

func (s *JobService) Get(ctx context.Context) ([]domain.Job, error) {
	return s.repo.List(ctx)
}

func (s *JobService) Create(ctx context.Context, j *domain.Job) (*domain.Job, error) {
	return s.repo.Create(ctx, j)
}

func (s *JobService) Delete(ctx context.Context, jobID int64) error {
	return s.repo.Delete(ctx, jobID)
}

func (s *JobService) Start(ctx context.Context, jobID int64) (*domain.Job, error) {
	return s.repo.MarkRunning(ctx, jobID)
}

func (s *JobService) Finish(ctx context.Context, jobID int64, status domain.JobStatus, outputLog *string) (*domain.Job, error) {
	return s.repo.MarkFinished(ctx, jobID, status, outputLog)
}
