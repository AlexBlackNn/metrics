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
	GetAllMetrics(
		ctx context.Context,
	) ([]models.Metric, error)
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

	select {
	case <-ctx.Done():
		ms.log.Error("Deadline exceeded while updating metric", "metric", metric)
		return ctx.Err()
	default:
		log := ms.log.With(
			slog.String("info", "SERVICE LAYER: metrics_service.UpdateMetricValue"),
		)
		log.Info("starts update metric value")

		if metric.Type == "counter" {

			metricStorage, err := ms.metricsStorage.GetMetric(ctx, metric.Name)
			if errors.Is(err, memstorage.ErrMetricNotFound) {
				err = ms.metricsStorage.UpdateMetric(ctx, metric)
				if err != nil {
					ms.log.Error(err.Error())
					return ErrCouldNotUpdateMetric
				}
				return nil
			}

			if err != nil {
				ms.log.Error(err.Error())
				return ErrCouldNotUpdateMetric
			}

			metric.Value = metricStorage.Value.(uint64) + metric.Value.(uint64)
			err = ms.metricsStorage.UpdateMetric(ctx, metric)
			if err != nil {
				ms.log.Error(err.Error())
				return ErrCouldNotUpdateMetric
			}
			return nil
		}
		err := ms.metricsStorage.UpdateMetric(ctx, metric)
		if err != nil {
			ms.log.Error(err.Error())
			return ErrCouldNotUpdateMetric
		}
		return nil
	}
}

func (ms *MetricService) GetOneMetricValue(ctx context.Context, key string) (models.Metric, error) {
	select {
	case <-ctx.Done():
		ms.log.Error("Deadline exceeded while getting metric", "name", key)
		return models.Metric{}, ctx.Err()
	default:
		log := ms.log.With(
			slog.String("info", "SERVICE LAYER: metrics_service.GetOneMetricValue"),
		)
		log.Info("starts getting metric value")

		metric, err := ms.metricsStorage.GetMetric(ctx, key)
		if errors.Is(err, memstorage.ErrMetricNotFound) {
			return models.Metric{}, ErrMetricNotFound
		}
		if err != nil {
			return models.Metric{}, ErrCouldNotUpdateMetric
		}
		return metric, nil
	}
}

func (ms *MetricService) GetAllMetrics(ctx context.Context) ([]models.Metric, error) {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.GetAllMetrics"),
	)
	log.Info("starts getting all metrics")

	metrics, err := ms.metricsStorage.GetAllMetrics(ctx)
	if errors.Is(err, memstorage.ErrMetricNotFound) {
		return []models.Metric{}, ErrMetricNotFound
	}
	if err != nil {
		return []models.Metric{}, ErrCouldNotUpdateMetric
	}
	return metrics, nil
}
