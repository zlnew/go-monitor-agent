// Package snapshot
package snapshot

import (
	"sync"

	"horizonx-server/internal/domain"
)

type MetricsStore struct {
	mu   sync.RWMutex
	data map[int64]domain.Metrics
}

func NewMetricsStore() *MetricsStore {
	return &MetricsStore{
		data: make(map[int64]domain.Metrics),
	}
}

func (s *MetricsStore) Set(serverID int64, m domain.Metrics) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[serverID] = m
}

func (s *MetricsStore) Get(serverID int64) (domain.Metrics, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[serverID]
	return val, ok
}

func (s *MetricsStore) GetAll() map[int64]domain.Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[int64]domain.Metrics, len(s.data))
	for k, v := range s.data {
		result[k] = v
	}
	return result
}
