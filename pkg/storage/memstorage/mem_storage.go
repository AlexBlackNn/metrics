package memstorage

import (
	"context"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"sync"
)

type MemStorage struct {
	mutex sync.RWMutex
	db    map[string]models.MetricInteraction
}

// New inits mem storage (map structure)
func New() (*MemStorage, error) {
	return &MemStorage{db: make(map[string]models.MetricInteraction)}, nil
}

// UpdateMetric updates metric value in mem storage
func (s *MemStorage) UpdateMetric(
	ctx context.Context,
	metric models.MetricInteraction,
) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.db[metric.GetName()] = metric
	return nil
}

// GetMetric gets metric value from mem storage
func (s *MemStorage) GetMetric(
	ctx context.Context,
	name string,
) (models.MetricInteraction, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	metric, ok := s.db[name]
	if !ok {
		return &models.Metric[float64]{}, ErrMetricNotFound
	}
	return metric, nil
}

// GetAllMetrics gets metric value from mem storage
func (s *MemStorage) GetAllMetrics(
	ctx context.Context,
) ([]models.MetricInteraction, error) {
	var metrics []models.MetricInteraction
	if len(s.db) == 0 {
		return []models.MetricInteraction{}, ErrMetricNotFound
	}
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, oneMetric := range s.db {
		metrics = append(metrics, oneMetric)
	}
	return metrics, nil
}
