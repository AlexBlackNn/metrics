package memstorage

import (
	"context"
	"errors"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"log/slog"
	"sync"
	"time"
)

var ErrFailedToRestoreMetrics = errors.New("failed to restore metrics")

type StateManager interface {
	saveMetrics() error
	restoreMetrics() error
}

type MemStorage struct {
	mutex    *sync.RWMutex
	db       dataBase
	cfg      *configserver.Config
	sm       StateManager
	log      *slog.Logger
	saveChan chan struct{}
}

// New inits mem storage (map structure)
func New(cfg *configserver.Config, log *slog.Logger) (*MemStorage, error) {
	db := make(dataBase)
	mutex := &sync.RWMutex{}
	memStorage := MemStorage{
		mutex:    mutex,
		cfg:      cfg,
		db:       db,
		log:      log,
		sm:       &dataBaseGOBStateManager{cfg: cfg, db: db, mutex: mutex, log: log},
		saveChan: make(chan struct{}),
	}

	go func() {
		memStorage.saveMetricToDisk()
	}()

	if cfg.ServerRestore {
		_ = memStorage.sm.restoreMetrics()
	}
	return &memStorage, nil
}

// saveMetricToDisk saves metrics to disk.
func (ms *MemStorage) saveMetricToDisk() {
	log := ms.log.With(
		slog.String("info", "STORAGE LAYER: mem_storage.saveMetricToDisk"),
	)
	log.Info("starts saving metric to disk")
	storeInterval := time.Duration(ms.cfg.ServerStoreInterval) * time.Second
	for {
		if ms.cfg.ServerStoreInterval > 0 {
			<-time.After(storeInterval)
			err := ms.sm.saveMetrics()
			if err != nil {
				log.Error("failed save metrics", "err", err)
			}
		} else {
			<-ms.saveChan
			err := ms.sm.saveMetrics()
			if err != nil {
				log.Error("failed save metrics", "err", err)
			}
		}
	}
}

// UpdateMetric updates metric value in mem storage.
func (ms *MemStorage) UpdateMetric(
	ctx context.Context,
	metric models.MetricGetter,
) error {
	ms.mutex.Lock()
	ms.db[metric.GetName()] = metric
	ms.mutex.Unlock()
	if ms.cfg.ServerStoreInterval == 0 {
		ms.saveChan <- struct{}{}
	}
	return nil
}

// GetMetric gets metric value from mem storage.
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

// GetAllMetrics gets metric value from mem storage.
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
