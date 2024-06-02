package metrics_service

import (
	"context"
	"errors"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/storage/mem_storage"
	"log/slog"
	"strconv"
	"strings"
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

func (ms *MetricService) UpdateMetricValue(ctx context.Context, urlPath string) error {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.UpdateMetricValue"),
	)
	log.Info("starts update metric value")

	parts := strings.Split(strings.Trim(urlPath, "/"), "/")
	if len(parts) != 4 {
		return ErrNotValidUrl
	}

	metricType := parts[1]
	metricName := parts[2]

	var err error
	var metricValue any
	var metric models.Metric

	switch metricType {
	case "gauge":
		metricValue, err = strconv.ParseFloat(parts[3], 64)
		if err != nil {
			return ErrNotValidMetricValue
		}
		metric = models.Metric{Type: metricType, Name: metricName, Value: metricValue}

	case "counter":
		metricValue, err = strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return ErrNotValidMetricValue
		}

		// Get existing metric from storage
		metric, err = ms.metricsStorage.GetMetric(ctx, metricName)
		if errors.Is(err, mem_storage.ErrMetricNotFound) {
			metric = models.Metric{Type: metricType, Name: metricName, Value: metricValue}
		} else {
			metric.Value = metric.Value.(int64) + metricValue.(int64)
		}
	default:
		return ErrNotValidMetricType
	}
	err = ms.metricsStorage.UpdateMetric(ctx, metric)
	if err != nil {
		return ErrCouldNotUpdateMetric
	}
	return nil
}
