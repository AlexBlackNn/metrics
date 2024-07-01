package memstorage

import (
	"context"
	"errors"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"sync"
)

var ErrFailedToRestoreMetrics = errors.New("failed to restore metrics")

type MemStorage struct {
	mutex sync.RWMutex
	db    dataBase
	cfg   *config.Config
	jm    *dataBaseJSONStateManager
}

// New inits mem storage (map structure)
func New(cfg *config.Config) (*MemStorage, error) {
	db := make(dataBase)
	memStorage := MemStorage{
		mutex: sync.RWMutex{},
		cfg:   cfg,
		db:    db,
		jm:    &dataBaseJSONStateManager{cfg: cfg, db: db},
	}

	if cfg.ServerRestore {
		err := memStorage.jm.restoreMetrics()
		if err != nil {
			if errors.Is(err, ErrFailedToRestoreMetrics) {
				return &memStorage, nil
			}
			return &memStorage, nil
		}
		return &memStorage, nil
	}
	return &memStorage, nil
}

// UpdateMetric updates metric value in mem storage
func (ms *MemStorage) UpdateMetric(
	ctx context.Context,
	metric models.MetricGetter,
) error {
	ms.mutex.Lock()
	ms.db[metric.GetName()] = metric
	ms.mutex.Unlock()
	go func() {
		ms.mutex.Lock()
		_ = ms.jm.saveMetrics()
		ms.mutex.Unlock()
	}()
	return nil
}

// GetMetric gets metric value from mem storage
func (ms *MemStorage) GetMetric(
	ctx context.Context,
	name string,
) (models.MetricGetter, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	metric, ok := ms.db[name]
	if !ok {
		return &models.Metric[float64]{}, ErrMetricNotFound
	}
	return metric, nil
}

// GetAllMetrics gets metric value from mem storage
func (ms *MemStorage) GetAllMetrics(
	ctx context.Context,
) ([]models.MetricGetter, error) {
	var metrics []models.MetricGetter
	if len(ms.db) == 0 {
		return []models.MetricGetter{}, ErrMetricNotFound
	}
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	for _, oneMetric := range ms.db {
		metrics = append(metrics, oneMetric)
	}
	return metrics, nil
}
