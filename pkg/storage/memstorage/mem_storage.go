package memstorage

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"io"
	"os"
	"sync"
)

var ErrFailedToRestoreMetrics = errors.New("failed to restore metrics")

type MemStorage struct {
	mutex sync.RWMutex
	db    DataBase
	cfg   *config.Config
}

// New inits mem storage (map structure)
func New(cfg *config.Config) (*MemStorage, error) {

	memStorage := MemStorage{
		mutex: sync.RWMutex{},
		cfg:   cfg,
		db:    make(DataBase),
	}
	if cfg.ServerRestore {
		err := memStorage.RestoreMetrics()
		if err != nil {
			if errors.Is(err, ErrFailedToRestoreMetrics) {
				return &memStorage, nil
			}
			return nil, err
		}
		return &memStorage, nil
	}
	return &memStorage, nil
}

func (s *MemStorage) RestoreMetrics() error {
	fmt.Println("START RESTORE METRICS")
	file, err := os.OpenFile(s.cfg.ServerFileStoragePath, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	reader := bufio.NewReader(file)

	tmpBuffer, err := io.ReadAll(reader)
	if err != nil {
		return ErrFailedToRestoreMetrics
	}
	err = s.db.decode(tmpBuffer)
	if err != nil {
		return ErrFailedToRestoreMetrics
	}
	return nil
}

func (s *MemStorage) SaveMetrics() error {

	file, err := os.OpenFile(s.cfg.ServerFileStoragePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	dataBaseBytes, err := s.db.encode()
	if err != nil {
		return err
	}
	_, err = writer.Write(dataBaseBytes)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMetric updates metric value in mem storage
func (s *MemStorage) UpdateMetric(
	ctx context.Context,
	metric models.MetricInteraction,
) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.db[metric.GetName()] = metric
	err := s.SaveMetrics()
	if err != nil {
		return err
	}
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
