package metrics

import (
	"context"
	"fmt"
	"log"
	"time"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/storage/snapshot"
	"horizonx-server/internal/transport/websocket"
)

type Service struct {
	repo      domain.MetricsRepository
	hub       *websocket.Hub
	snapshot  *snapshot.MetricsStore
	saveQueue chan domain.Metrics
}

func NewService(repo domain.MetricsRepository, snapshot *snapshot.MetricsStore, hub *websocket.Hub) *Service {
	s := &Service{
		repo:      repo,
		hub:       hub,
		snapshot:  snapshot,
		saveQueue: make(chan domain.Metrics, 1000),
	}

	go s.worker()
	return s
}

func (s *Service) Ingest(ctx context.Context, m domain.Metrics) error {
	s.snapshot.Set(m.ServerID, m)

	s.hub.Emit(fmt.Sprintf("server:%d:metrics", m.ServerID), "metrics.updated", m)

	select {
	case s.saveQueue <- m:
	default:
		log.Println("WARNING: Metric queue full! Dropping data.")
	}

	return nil
}

func (s *Service) worker() {
	buffer := make([]domain.Metrics, 0, 100)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	flush := func() {
		if len(buffer) > 0 {
			if err := s.repo.BulkInsert(context.Background(), buffer); err != nil {
				log.Printf("Failed to flush metrics to DB: %v", err)
			}
			buffer = buffer[:0]
		}
	}

	for {
		select {
		case m := <-s.saveQueue:
			buffer = append(buffer, m)
			if len(buffer) >= 100 {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}
