package metricsservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/pkg/storage"
	"log/slog"
)

type MetricsStorage interface {
	UpdateMetric(
		ctx context.Context,
		metric models.MetricGetter,
	) error
	UpdateSeveralMetrics(
		ctx context.Context,
		metrics []models.MetricGetter,
	) error
	GetMetric(
		ctx context.Context,
		metric models.MetricGetter,
	) (models.MetricGetter, error)
	GetAllMetrics(
		ctx context.Context,
	) ([]models.MetricGetter, error)
}

type HealthChecker interface {
	HealthCheck(
		ctx context.Context,
	) error
}

type MetricService struct {
	log            *slog.Logger
	cfg            *configserver.Config
	metricsStorage MetricsStorage
	healthChecker  HealthChecker
}

// New returns a new instance of MonitoringService.
func New(
	log *slog.Logger,
	cfg *configserver.Config,
	metricsStorage MetricsStorage,
	healthChecker HealthChecker,
) *MetricService {
	return &MetricService{
		log:            log,
		cfg:            cfg,
		metricsStorage: metricsStorage,
		healthChecker:  healthChecker,
	}
}

// UpdateMetricValue updates metric value or create new metric.
func (ms *MetricService) UpdateMetricValue(ctx context.Context, metric models.MetricInteraction) error {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.UpdateMetricValue"),
	)
	log.Info("starts update metric value")

	if metric.GetType() == configserver.MetricTypeCounter {

		metricStorage, err := ms.metricsStorage.GetMetric(ctx, metric)
		if errors.Is(err, storage.ErrMetricNotFound) {
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

// UpdateSeveralMetrics updates several metric value or create new ones.
func (ms *MetricService) UpdateSeveralMetrics(ctx context.Context, metrics []models.MetricInteraction) error {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.UpdateSeveralMetrics"),
	)
	log.Info("starts update several metric values")

	for _, OneMetric := range metrics {
		fmt.Println("4444444444444", OneMetric, OneMetric.GetType(), OneMetric.GetName(), OneMetric.GetValue())
	}

	var errs []error
	for _, oneMetric := range metrics {
		if oneMetric.GetType() == configserver.MetricTypeCounter {

			metricStorage, err := ms.metricsStorage.GetMetric(ctx, oneMetric)
			if errors.Is(err, storage.ErrMetricNotFound) {
				fmt.Println("4444444444444*", ErrMetricNotFound, oneMetric.GetType(), oneMetric.GetName(), oneMetric.GetValue())
				err = ms.metricsStorage.UpdateMetric(ctx, oneMetric)
				if err != nil {
					ms.log.Error(err.Error())
					errs = append(errs, ErrCouldNotUpdateMetric)
					continue
				}
				continue
			}
			if err != nil {
				ms.log.Error(err.Error())
				errs = append(errs, ErrCouldNotUpdateMetric)
				continue
			}
			fmt.Println("4444444444444**", oneMetric.GetType(), oneMetric.GetName(), oneMetric.GetValue())
			fmt.Println("4444444444444***", metricStorage.GetType(), metricStorage.GetName(), metricStorage.GetValue())
			err = oneMetric.AddValue(metricStorage)
			fmt.Println("4444444444444****", ErrMetricNotFound, oneMetric.GetType(), oneMetric.GetName(), oneMetric.GetValue())
			if err != nil {
				ms.log.Error(err.Error())
				errs = append(errs, ErrCouldNotUpdateMetric)
				continue
			}
			err = ms.metricsStorage.UpdateMetric(ctx, oneMetric)
			if err != nil {
				ms.log.Error(err.Error())
				errs = append(errs, ErrCouldNotUpdateMetric)
				continue
			}
			continue
		}
		err := ms.metricsStorage.UpdateMetric(ctx, oneMetric)
		if err != nil {
			ms.log.Error(err.Error())
			errs = append(errs, ErrCouldNotUpdateMetric)
			continue
		}
		log.Info("finish updating metric value")
	}
	return errors.Join(errs...)
}

// GetOneMetricValue extracts metric.
func (ms *MetricService) GetOneMetricValue(ctx context.Context, metric models.MetricGetter) (models.MetricGetter, error) {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.GetOneMetricValue"),
	)
	log.Info("starts getting metric value")
	metric, err := ms.metricsStorage.GetMetric(ctx, metric)
	if err != nil {
		if errors.Is(err, storage.ErrMetricNotFound) {
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
	if errors.Is(err, storage.ErrMetricNotFound) {
		return nil, ErrMetricNotFound
	}
	if err != nil {
		return nil, ErrCouldNotUpdateMetric
	}
	log.Info("finish getting all metrics")
	return metrics, nil
}

// HealthCheck returns service health check
func (ms *MetricService) HealthCheck(ctx context.Context) error {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.HealthCheck"),
	)
	log.Info("starts getting health check")
	defer log.Info("finish getting health check")
	return ms.healthChecker.HealthCheck(ctx)
}
