package metricsservice

import (
	"context"
	"errors"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/storage/memstorage"
	"log/slog"
)

type MetricsStorageInterface interface {
	UpdateMetric(
		ctx context.Context,
		metric models.Metric,
	) error
	GetMetric(
		ctx context.Context,
		metricName string,
	) (models.Metric, error)
}

type MetricService struct {
	log            *slog.Logger
	cfg            *config.Config
	metricsStorage MetricsStorageInterface
}

// New returns a new instance of MonitoringService
func New(
	log *slog.Logger,
	cfg *config.Config,
	metricsStorage MetricsStorageInterface,
) *MetricService {
	return &MetricService{
		log:            log,
		cfg:            cfg,
		metricsStorage: metricsStorage,
	}
}

func (ms *MetricService) UpdateMetricValue(ctx context.Context, metric models.Metric) error {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.UpdateMetricValue"),
	)
	log.Info("starts update metric value")

	if metric.Type == "counter" {

		// Get existing metric from storage
		metricStorage, err := ms.metricsStorage.GetMetric(ctx, metric.Name)
		if !errors.Is(err, memstorage.ErrMetricNotFound) {
			metric.Value = metricStorage.Value.(uint64) + metric.Value.(uint64)
		}
		err = ms.metricsStorage.UpdateMetric(ctx, metric)
		if err != nil {
			return ErrCouldNotUpdateMetric
		}
		return nil
	}
	err := ms.metricsStorage.UpdateMetric(ctx, metric)
	if err != nil {
		return ErrCouldNotUpdateMetric
	}
	return nil
}
