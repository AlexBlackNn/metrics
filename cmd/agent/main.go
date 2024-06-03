package main

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/cmd/agent/domain/models"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type MetricsMonitor struct {
	Metrics      map[string]models.Metric
	PollInterval int
	mutex        sync.RWMutex
}

func NewMetricsMonitor(pollInterval int) *MetricsMonitor {
	return &MetricsMonitor{
		Metrics:      make(map[string]models.Metric),
		PollInterval: pollInterval,
	}
}

func (mm *MetricsMonitor) Start() {
	var rtm runtime.MemStats
	interval := time.Duration(mm.PollInterval) * time.Second

	mm.Metrics["PollCount"] = models.Metric{Type: "counter", Value: int64(0), Name: "PollCount"}
	for {
		// Read full mem stats
		runtime.ReadMemStats(&rtm)
		t := reflect.TypeOf(rtm)
		if t.Kind() == reflect.Struct {
			mm.mutex.Lock()
			for i := 0; i < t.NumField(); i++ {
				metricName := t.Field(i).Name
				metricValue := reflect.ValueOf(rtm).FieldByName(metricName).Interface()
				metricType := reflect.TypeOf(metricValue).String()
				if metricType == "float64" || metricType == "uint32" || metricType == "uint64" {
					mm.Metrics[metricName] = models.Metric{Type: "gauge", Value: metricValue, Name: metricName}
				}
			}
			mm.Metrics["PollCount"] = models.Metric{Type: "counter", Value: mm.Metrics["PollCount"].Value.(int64) + 1, Name: "PollCount"}
			mm.Metrics["RandomValue"] = models.Metric{Type: "gauge", Value: rand.Int63(), Name: "RandomValue"}
			mm.mutex.Unlock()
			fmt.Println(mm.Metrics)
			<-time.After(interval)
		}
	}
}

//func main() {
//	metricsMonitor := NewMetricsMonitor(10000)
//	metricsMonitor.Start()
//}

func main() {
	var wg sync.WaitGroup
	metricsMonitor := NewMetricsMonitor(10000)

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsMonitor.Start()
	}()
	wg.Wait()
}
