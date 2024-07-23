package agentmetricsservice

import (
	"context"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
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

// Collect starts collecting runtime metrics.
func (ms *MonitorService) Collect(ctx context.Context) {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: agentmetricservice.Start"),
	)

	var rtm runtime.MemStats
	ms.Metrics["PollCount"] = &models.Metric[uint64]{Type: configagent.MetricTypeCounter, Value: uint64(0), Name: "PollCount"}
	for {
		select {
		case <-ctx.Done():
			log.Info("stop signal received")
			return
		case <-time.After(time.Duration(ms.cfg.PollInterval) * time.Second):
			log.Info("starts Collect metric pooling")
			// Read full mem stats
			runtime.ReadMemStats(&rtm)
			ms.mutex.Lock()
			ms.Metrics["Alloc"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.Alloc, Name: "Alloc"}
			ms.Metrics["BuckHashSys"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.BuckHashSys, Name: "BuckHashSys"}
			ms.Metrics["Frees"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.Frees, Name: "Frees"}
			ms.Metrics["GCCPUFraction"] = &models.Metric[float64]{Type: configagent.MetricTypeGauge, Value: rtm.GCCPUFraction, Name: "GCCPUFraction"}
			ms.Metrics["GCSys"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.GCSys, Name: "GCSys"}
			ms.Metrics["HeapAlloc"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.HeapAlloc, Name: "HeapAlloc"}
			ms.Metrics["HeapIdle"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.HeapIdle, Name: "HeapIdle"}
			ms.Metrics["HeapInuse"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.HeapInuse, Name: "HeapInuse"}
			ms.Metrics["HeapObjects"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.HeapObjects, Name: "HeapObjects"}
			ms.Metrics["HeapReleased"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.HeapReleased, Name: "HeapReleased"}
			ms.Metrics["HeapSys"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.HeapSys, Name: "HeapSys"}
			ms.Metrics["LastGC"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.LastGC, Name: "LastGC"}
			ms.Metrics["Lookups"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.Lookups, Name: "Lookups"}
			ms.Metrics["MCacheInuse"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.MCacheInuse, Name: "MCacheInuse"}
			ms.Metrics["MCacheSys"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.MCacheSys, Name: "MCacheSys"}
			ms.Metrics["MSpanInuse"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.MSpanInuse, Name: "MSpanInuse"}
			ms.Metrics["MSpanSys"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.MSpanSys, Name: "MSpanSys"}
			ms.Metrics["Mallocs"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.Mallocs, Name: "Mallocs"}
			ms.Metrics["NextGC"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.NextGC, Name: "NextGC"}
			ms.Metrics["NumForcedGC"] = &models.Metric[uint32]{Type: configagent.MetricTypeGauge, Value: rtm.NumForcedGC, Name: "NumForcedGC"}
			ms.Metrics["NumGC"] = &models.Metric[uint32]{Type: configagent.MetricTypeGauge, Value: rtm.NumGC, Name: "NumGC"}
			ms.Metrics["OtherSys"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.OtherSys, Name: "OtherSys"}
			ms.Metrics["PauseTotalNs"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.PauseTotalNs, Name: "PauseTotalNs"}
			ms.Metrics["StackInuse"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.StackInuse, Name: "StackInuse"}
			ms.Metrics["StackSys"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.StackSys, Name: "StackSys"}
			ms.Metrics["Sys"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.Sys, Name: "Sys"}
			ms.Metrics["TotalAlloc"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rtm.TotalAlloc, Name: "TotalAlloc"}
			ms.Metrics["PollCount"] = &models.Metric[uint64]{Type: configagent.MetricTypeCounter, Value: ms.Metrics["PollCount"].GetValue().(uint64) + 1, Name: "PollCount"}
			ms.Metrics["RandomValue"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: rand.Uint64(), Name: "RandomValue"}
			ms.mutex.Unlock()
			log.Info("metric pooling finished")
		}
	}
}

// CollectAddition Collect starts collecting gopsutil metrics.
func (ms *MonitorService) CollectAddition(ctx context.Context) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: agentmetricservice.Start"),
	)

	virtMem, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err.Error())
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("stop signal received")
			return
		case <-time.After(time.Duration(ms.cfg.PollInterval) * time.Second):
			log.Info("starts CollectAddingMetrics metrics pooling")

			utilCPU, err := ms.calculateUtilization()
			if err != nil {
				log.Error(err.Error())
			}
			ms.Metrics["CPUutilization1"] = &models.Metric[float64]{Type: configagent.MetricTypeGauge, Value: utilCPU, Name: "CPUutilization1"}
			ms.Metrics["TotalMemory"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: virtMem.Total, Name: "TotalMemory"}
			ms.Metrics["FreeMemory"] = &models.Metric[uint64]{Type: configagent.MetricTypeGauge, Value: virtMem.Available, Name: "FreeMemory"}
			log.Info("metric pooling finished")
		}
	}
}

// GetMetrics return collected metrics as thread safe map.
func (ms *MonitorService) GetMetrics() map[string]models.MetricInteraction {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	return ms.Metrics
}

func (ms *MonitorService) calculateUtilization() (float64, error) {
	// get available cpu
	numCPUs := runtime.NumCPU()

	// Get cpu loading statistic
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}

	// calculate average CPU loading
	totalPercent := 0.0
	for _, percent := range cpuPercent {
		totalPercent += percent
	}
	return totalPercent / float64(numCPUs), nil
}
