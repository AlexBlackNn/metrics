package agentmetricsservice

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"log/slog"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type MetricsService struct {
	log     *slog.Logger
	cfg     *config.Config
	Metrics map[string]models.Metric
	mutex   sync.RWMutex
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *MetricsService {
	return &MetricsService{
		Metrics: make(map[string]models.Metric),
		log:     log,
		cfg:     cfg,
	}
}

func (ms *MetricsService) Start(stop chan struct{}) {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: metrics_service.UpdateMetricValue"),
	)
	log.Info("starts update metric value")

	var rtm runtime.MemStats
	interval := time.Duration(ms.cfg.PollInterval) * time.Second

	ms.Metrics["PollCount"] = models.Metric{Type: "counter", Value: int64(0), Name: "PollCount"}
	for {
		select {
		case <-stop:
			return
		default:
			// Read full mem stats
			runtime.ReadMemStats(&rtm)
			t := reflect.TypeOf(rtm)
			if t.Kind() == reflect.Struct {
				ms.mutex.Lock()
				for i := 0; i < t.NumField(); i++ {
					metricName := t.Field(i).Name
					metricValue := reflect.ValueOf(rtm).FieldByName(metricName).Interface()
					metricType := reflect.TypeOf(metricValue).String()
					if metricType == "float64" || metricType == "uint32" || metricType == "uint64" {
						ms.Metrics[metricName] = models.Metric{Type: "gauge", Value: metricValue, Name: metricName}
					}
				}
				ms.Metrics["PollCount"] = models.Metric{Type: "counter", Value: ms.Metrics["PollCount"].Value.(int64) + 1, Name: "PollCount"}
				ms.Metrics["RandomValue"] = models.Metric{Type: "gauge", Value: rand.Int63(), Name: "RandomValue"}
				ms.mutex.Unlock()
				<-time.After(interval)
			}
		}
	}
}

func (ms *MetricsService) GetMetrics() map[string]models.Metric {
	return ms.Metrics
}
