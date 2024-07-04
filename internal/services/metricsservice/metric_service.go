package metricsservice

import (
	"context"
	"errors"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/pkg/storage/memstorage"
	"log/slog"
)

type MetricsStorage interface {
	UpdateMetric(
		ctx context.Context,
		metric models.MetricGetter,
	) error
	GetMetric(
		ctx context.Context,
		metricName string,
	) (models.MetricGetter, error)
	GetAllMetrics(
		ctx context.Context,
	) ([]models.MetricGetter, error)
}

type MetricService struct {
	log            *slog.Logger
	cfg            *configserver.Config
	metricsStorage MetricsStorage
}

// New returns a new instance of MonitoringService.
func New(
	log *slog.Logger,
	cfg *configserver.Config,
	metricsStorage MetricsStorage,
) *MetricService {
	return &MetricService{
		log:            log,
		cfg:            cfg,
		metricsStorage: metricsStorage,
	}
}

// UpdateMetricValue updates metric value or create new metric.
func (ms *MetricService) UpdateMetricValue(ctx context.Context, metric models.MetricInteraction) error {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.UpdateMetricValue"),
	)
	log.Info("starts update metric value")

	if metric.GetType() == configserver.MetricTypeCounter {

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

// GetOneMetricValue extracts metric.
func (ms *MetricService) GetOneMetricValue(ctx context.Context, key string) (models.MetricGetter, error) {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.GetOneMetricValue"),
	)
	log.Info("starts getting metric value")
	metric, err := ms.metricsStorage.GetMetric(ctx, key)
	if err != nil {
		if errors.Is(err, memstorage.ErrMetricNotFound) {
			log.Warn("metric not found")
			return nil, ErrMetricNotFound
		}
		log.Error(err.Error())
		return nil, ErrCouldNotGetMetric
	}
	log.Info("finish getting metric value")
	return metric, nil
}

// GetAllMetrics extracts all metric.
func (ms *MetricService) GetAllMetrics(ctx context.Context) ([]models.MetricGetter, error) {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.GetAllMetrics"),
	)
	log.Info("starts getting all metrics")

	metrics, err := ms.metricsStorage.GetAllMetrics(ctx)
	if errors.Is(err, memstorage.ErrMetricNotFound) {
		return nil, ErrMetricNotFound
	}
	if err != nil {
		return nil, ErrCouldNotUpdateMetric
	}
	log.Info("finish getting all metrics")
	return metrics, nil
}
