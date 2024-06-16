package metricsservice

import (
	"context"
	"errors"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/pkg/storage/memstorage"
	"log/slog"
)

type MetricsStorage interface {
	UpdateMetric(
		ctx context.Context,
		metric models.MetricInteraction,
	) error
	GetMetric(
		ctx context.Context,
		metricName string,
	) (models.MetricInteraction, error)
	GetAllMetrics(
		ctx context.Context,
	) ([]models.MetricInteraction, error)
}

type MetricService struct {
	log            *slog.Logger
	cfg            *config.Config
	metricsStorage MetricsStorage
}

// New returns a new instance of MonitoringService
func New(
	log *slog.Logger,
	cfg *config.Config,
	metricsStorage MetricsStorage,
) *MetricService {
	return &MetricService{
		log:            log,
		cfg:            cfg,
		metricsStorage: metricsStorage,
	}
}

func (ms *MetricService) UpdateMetricValue(ctx context.Context, metric models.MetricInteraction) error {

	select {
	case <-ctx.Done():
		ms.log.Error("Deadline exceeded while updating metric", "metric", metric)
		return ctx.Err()
	default:
		log := ms.log.With(
			slog.String("info", "SERVICE LAYER: metrics_service.UpdateMetricValue"),
		)
		log.Info("starts update metric value")

		if metric.GetType() == "counter" {

			metricStorage, err := ms.metricsStorage.GetMetric(ctx, metric.GetName())
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
			err = metric.AddValue(metricStorage)
			if err != nil {
				ms.log.Error(err.Error())
				return ErrCouldNotUpdateMetric
			}
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
		log.Info("finish updating metric value")
		return nil
	}
}

func (ms *MetricService) GetOneMetricValue(ctx context.Context, key string) (models.MetricInteraction, error) {
	select {
	case <-ctx.Done():
		ms.log.Error("Deadline exceeded while getting metric", "name", key)
		return &models.Metric[float64]{}, ctx.Err()
	default:
		log := ms.log.With(
			slog.String("info", "SERVICE LAYER: metrics_service.GetOneMetricValue"),
		)
		log.Info("starts getting metric value")

		metric, err := ms.metricsStorage.GetMetric(ctx, key)
		if errors.Is(err, memstorage.ErrMetricNotFound) {
			return &models.Metric[float64]{}, ErrMetricNotFound
		}
		if err != nil {
			return &models.Metric[float64]{}, ErrCouldNotUpdateMetric
		}
		log.Info("finish getting metric value")
		return metric, nil
	}
}

func (ms *MetricService) GetAllMetrics(ctx context.Context) ([]models.MetricInteraction, error) {
	select {
	case <-ctx.Done():
		ms.log.Error("Deadline exceeded while getting all metrics")
		return []models.MetricInteraction{}, ctx.Err()
	default:
		log := ms.log.With(
			slog.String("info", "SERVICE LAYER: metrics_service.GetAllMetrics"),
		)
		log.Info("starts getting all metrics")

		metrics, err := ms.metricsStorage.GetAllMetrics(ctx)
		if errors.Is(err, memstorage.ErrMetricNotFound) {
			return []models.MetricInteraction{}, ErrMetricNotFound
		}
		if err != nil {
			return []models.MetricInteraction{}, ErrCouldNotUpdateMetric
		}
		log.Info("finish getting all metrics")
		return metrics, nil
	}
}
