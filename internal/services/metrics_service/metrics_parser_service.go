package metrics_service

import (
	"context"
	"fmt"
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

func (ms *MetricService) UpdateMetricValue(
	ctx context.Context,
	urlPath string,
) error {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.UpdateMetricValue"),
	)
	log.Info("starts update metric value")

	//
	parts := strings.Split(strings.Trim(urlPath, "/"), "/")
	if len(parts) != 4 {
		return ErrNotValidUrl
	}

	var (
		metricType  string
		metricValue any
		err         error
		metric      models.Metric
	)

	switch metricType = parts[1]; metricType {
	case "gauge":
		metricValue, err = strconv.ParseFloat(parts[3], 64)
		if err != nil {
			return ErrNotValidMetricValue
		}
	case "counter":
		metricValue, err = strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return ErrNotValidMetricValue
		}
	default:
		return ErrNotValidMetricType
	}

	metric = models.Metric{Type: metricType, Name: parts[2], Value: metricValue}
	err = ms.metricsStorage.UpdateMetric(ctx, metric)
	gotMetric, err := ms.metricsStorage.GetMetric(ctx, metric.Name)
	fmt.Println("111111111111", gotMetric)
	if err != nil {
		return fmt.Errorf("could not update metric: %w", err)
	}
	return nil
}
