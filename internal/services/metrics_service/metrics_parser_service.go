package metrics_service

import (
	"context"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config"
	"log/slog"
	"strconv"
	"strings"
)

type MetricService struct {
	log *slog.Logger
	cfg *config.Config
}

// New returns a new instance of MonitoringService
func New(
	log *slog.Logger,
	cfg *config.Config,
) *MetricService {
	return &MetricService{
		log: log,
		cfg: cfg,
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
	fmt.Println(metricType, metricValue)

	return nil
}
