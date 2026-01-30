package redis

import (
	"context"

	"horizonx/internal/domain"

	"github.com/redis/go-redis/v9"
)

const maxLen = 5000

type MetricsRegistry struct {
	registry *Registry
	key      string
}

func NewMetricsRegistry(r *redis.Client, key string) *MetricsRegistry {
	return &MetricsRegistry{
		registry: NewRegistry(r),
		key:      key,
	}
}

func (r *MetricsRegistry) Append(ctx context.Context, m *domain.Metrics) (string, error) {
	return r.registry.Append(ctx, r.key, m, maxLen)
}

func (r *MetricsRegistry) GetBatch(ctx context.Context, limit int64) ([]domain.Metrics, []string, error) {
	messages, err := r.registry.GetRange(ctx, r.key, limit)
	if err != nil {
		return nil, nil, err
	}

	items, ids, err := ParseStreamMessages[domain.Metrics](messages)
	if err != nil {
		return nil, nil, err
	}

	return items, ids, nil
}

func (r *MetricsRegistry) GetLatest(ctx context.Context) (*domain.Metrics, string, error) {
	messages, err := r.registry.GetLatest(ctx, r.key)
	if err != nil {
		return nil, "", err
	}

	if len(messages) == 0 {
		return nil, "", nil
	}

	items, ids, err := ParseStreamMessages[domain.Metrics](messages)
	if err != nil {
		return nil, "", err
	}

	if len(items) == 0 {
		return nil, "", nil
	}

	return &items[0], ids[0], nil
}

func (r *MetricsRegistry) Ack(ctx context.Context, ids []string) error {
	return r.registry.Ack(ctx, r.key, ids)
}
