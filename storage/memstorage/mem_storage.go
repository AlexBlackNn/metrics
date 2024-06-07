package memstorage

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
	metric, ok := s.db[name]
	if !ok {
		return models.Metric{}, ErrMetricNotFound
	}
	return metric, nil
}

func (s *Storage) GetAllMetrics(
	ctx context.Context,
) ([]models.Metric, error) {
	var metrics []models.Metric

	if len(s.db) == 0 {
		return []models.Metric{}, ErrMetricNotFound
	}
	for _, oneMetric := range s.db {
		metrics = append(metrics, oneMetric)
	}

	return metrics, nil
}
