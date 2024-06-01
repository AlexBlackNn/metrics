package mem_storage

import (
	"context"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
)

type Storage struct {
	db map[string]models.Metric
}

func New() (*Storage, error) {
	return &Storage{db: make(map[string]models.Metric)}, nil
}

func (s *Storage) UpdateMetric(
	ctx context.Context,
	metric models.Metric,
) error {
	s.db[metric.Name] = metric
	return nil
}

func (s *Storage) GetMetric(
	ctx context.Context,
	name string,
) (models.Metric, error) {
	return s.db[name], nil
}
