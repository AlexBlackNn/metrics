package agentmetricsservice

import (
	"context"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"log/slog"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type MonitorService struct {
	log     *slog.Logger
	cfg     *configagent.Config
	Metrics map[string]models.MetricInteraction
	mutex   sync.RWMutex
}

func New(
	log *slog.Logger,
	cfg *configagent.Config,
) *MonitorService {
	return &MonitorService{
		Metrics: make(map[string]models.MetricInteraction),
		log:     log,
		cfg:     cfg,
	}
}

// Collect starts collecting runtime metrics
func (ms *MonitorService) Collect(ctx context.Context) {
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
			ms.mutex.Lock()
			ms.Metrics["Alloc"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.Alloc, Name: "Alloc"}
			ms.Metrics["BuckHashSys"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.BuckHashSys, Name: "BuckHashSys"}
			ms.Metrics["Frees"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.Frees, Name: "Frees"}
			ms.Metrics["GCCPUFraction"] = &models.Metric[float64]{Type: "gauge", Value: rtm.GCCPUFraction, Name: "GCCPUFraction"}
			ms.Metrics["GCSys"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.GCSys, Name: "GCSys"}
			ms.Metrics["HeapAlloc"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.HeapAlloc, Name: "HeapAlloc"}
			ms.Metrics["HeapIdle"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.HeapIdle, Name: "HeapIdle"}
			ms.Metrics["HeapInuse"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.HeapInuse, Name: "HeapInuse"}
			ms.Metrics["HeapObjects"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.HeapObjects, Name: "HeapObjects"}
			ms.Metrics["HeapReleased"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.HeapReleased, Name: "HeapReleased"}
			ms.Metrics["HeapSys"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.HeapSys, Name: "HeapSys"}
			ms.Metrics["LastGC"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.LastGC, Name: "LastGC"}
			ms.Metrics["Lookups"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.Lookups, Name: "Lookups"}
			ms.Metrics["MCacheInuse"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.MCacheInuse, Name: "MCacheInuse"}
			ms.Metrics["MCacheSys"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.MCacheSys, Name: "MCacheSys"}
			ms.Metrics["MSpanInuse"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.MSpanInuse, Name: "MSpanInuse"}
			ms.Metrics["MSpanSys"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.MSpanSys, Name: "MSpanSys"}
			ms.Metrics["Mallocs"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.Mallocs, Name: "Mallocs"}
			ms.Metrics["NextGC"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.NextGC, Name: "NextGC"}
			ms.Metrics["NumForcedGC"] = &models.Metric[uint32]{Type: "gauge", Value: rtm.NumForcedGC, Name: "NumForcedGC"}
			ms.Metrics["NumGC"] = &models.Metric[uint32]{Type: "gauge", Value: rtm.NumGC, Name: "NumGC"}
			ms.Metrics["OtherSys"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.OtherSys, Name: "OtherSys"}
			ms.Metrics["PauseTotalNs"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.PauseTotalNs, Name: "PauseTotalNs"}
			ms.Metrics["StackInuse"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.StackInuse, Name: "StackInuse"}
			ms.Metrics["StackSys"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.StackSys, Name: "StackSys"}
			ms.Metrics["Sys"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.Sys, Name: "Sys"}
			ms.Metrics["TotalAlloc"] = &models.Metric[uint64]{Type: "gauge", Value: rtm.TotalAlloc, Name: "TotalAlloc"}
			ms.Metrics["PollCount"] = &models.Metric[uint64]{Type: "counter", Value: ms.Metrics["PollCount"].GetValue().(uint64) + 1, Name: "PollCount"}
			ms.Metrics["RandomValue"] = &models.Metric[uint64]{Type: "gauge", Value: rand.Uint64(), Name: "RandomValue"}
			ms.mutex.Unlock()
			log.Info("metric pooling finished")
			<-time.After(time.Duration(ms.cfg.PollInterval) * time.Second)
		}
	}
}

// GetMetrics return collected metrics as thread safe map
func (ms *MonitorService) GetMetrics() map[string]models.MetricInteraction {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	return ms.Metrics
}
