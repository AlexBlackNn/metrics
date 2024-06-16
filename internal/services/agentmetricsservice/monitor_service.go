package agentmetricsservice

import (
	"context"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"log/slog"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type MonitorService struct {
	log     *slog.Logger
	cfg     *config.Config
	Metrics map[string]models.MetricInteraction
	mutex   sync.RWMutex
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *MonitorService {
	return &MonitorService{
		Metrics: make(map[string]models.MetricInteraction),
		log:     log,
		cfg:     cfg,
	}
}

// Start starts collecting runtime metrics
func (ms *MonitorService) Start(ctx context.Context) {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: agentmetricservice.Start"),
	)

	var rtm runtime.MemStats
	ms.Metrics["PollCount"] = &models.Metric[uint64]{Type: "counter", Value: uint64(0), Name: "PollCount"}
	for {
		select {
		case <-ctx.Done():
			log.Info("stop signal received")
			return
		default:
			log.Info("starts metric pooling")
			// Read full mem stats
			runtime.ReadMemStats(&rtm)
			t := reflect.TypeOf(rtm)
			if t.Kind() == reflect.Struct {
				ms.mutex.Lock()
				for i := 0; i < t.NumField(); i++ {
					metricName := t.Field(i).Name
					metricValue := reflect.ValueOf(rtm).FieldByName(metricName).Interface()
					metricType := reflect.TypeOf(metricValue).String()
					if metricType == "float64" {
						ms.Metrics[metricName] = &models.Metric[float64]{Type: "gauge", Value: metricValue.(float64), Name: metricName}
					}
					if metricType == "uint32" {
						ms.Metrics[metricName] = &models.Metric[uint32]{Type: "gauge", Value: metricValue.(uint32), Name: metricName}
					}
					if metricType == "uint64" {
						ms.Metrics[metricName] = &models.Metric[uint64]{Type: "gauge", Value: metricValue.(uint64), Name: metricName}
					}
				}
				ms.Metrics["PollCount"] = &models.Metric[uint64]{Type: "counter", Value: ms.Metrics["PollCount"].GetValue().(uint64) + 1, Name: "PollCount"}
				ms.Metrics["RandomValue"] = &models.Metric[uint64]{Type: "gauge", Value: rand.Uint64(), Name: "RandomValue"}
				ms.mutex.Unlock()
				log.Info("metric pooling finished")
				<-time.After(time.Duration(ms.cfg.PollInterval) * time.Second)
			}
		}
	}
}

// GetMetrics return collected metrics as thread safe map
func (ms *MonitorService) GetMetrics() map[string]models.MetricInteraction {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	return ms.Metrics
}
