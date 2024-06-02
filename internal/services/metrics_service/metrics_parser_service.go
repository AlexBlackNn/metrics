package metrics_service

import (
	"context"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
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

//func (ms *MetricService) UpdateMetricValue(
//	ctx context.Context,
//	urlPath string,
//) error {
//	log := ms.log.With(
//		slog.String("info", "SERVICE LAYER: metrics_service.UpdateMetricValue"),
//	)
//	log.Info("starts update metric value")
//
//	parts := strings.Split(strings.Trim(urlPath, "/"), "/")
//	if len(parts) != 4 {
//		return ErrNotValidUrl
//	}
//
//	var (
//		metricType  string
//		metricValue any
//		err         error
//		metric      models.Metric
//	)
//
//	switch metricType = parts[1]; metricType {
//	case "gauge":
//		metricValue, err = strconv.ParseFloat(parts[3], 64)
//		if err != nil {
//			return ErrNotValidMetricValue
//		}
//		metric = models.Metric{Type: metricType, Name: parts[2], Value: metricValue}
//
//	case "counter":
//		metricValue, err = strconv.ParseInt(parts[3], 10, 64)
//		if err != nil {
//			return ErrNotValidMetricValue
//		}
//		metric, err = ms.metricsStorage.GetMetric(ctx, metric.Name)
//		if err != nil {
//			return ErrCouldNotUpdateMetric
//		}
//		switch value := metric.Value.(type) {
//		case int64:
//			metric.Value = value + metricValue.(int64)
//		default:
//			metric.Value = metricValue.(int64)
//		}
//	default:
//		return ErrNotValidMetricType
//	}
//
//	err = ms.metricsStorage.UpdateMetric(ctx, metric)
//	if err != nil {
//		return ErrCouldNotUpdateMetric
//	}
//	return nil
//}

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
		if err != nil {
			return ErrCouldNotUpdateMetric
		}

		// Increment counter value
		switch value := metric.Value.(type) {
		case int64:
			metric.Value = value + metricValue.(int64)
		default:
			metric = models.Metric{Type: metricType, Name: metricName, Value: metricValue}
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
