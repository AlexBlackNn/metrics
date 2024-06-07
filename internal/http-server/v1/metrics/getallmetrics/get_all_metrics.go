package getallmetrics

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/appserver"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

//func New(log *slog.Logger, application *appserver.App) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		if r.Method != http.MethodGet {
//			w.WriteHeader(http.StatusMethodNotAllowed)
//			return
//		}
//
//		metrics, err := application.MetricsService.GetAllMetrics(context.Background())
//
//		if errors.Is(err, metricsservice.ErrMetricNotFound) {
//			http.Error(w, err.Error(), http.StatusNotFound)
//			return
//		}
//		fmt.Println(metrics)
//		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
//		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		w.WriteHeader(http.StatusOK)
//
//	}
//}
//

func New(log *slog.Logger, application *appserver.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		metrics, err := application.MetricsService.GetAllMetrics(context.Background())

		if errors.Is(err, metricsservice.ErrMetricNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Define Go template
		tmpl, err := template.New("metrics").ParseFiles("/home/alex/Dev/GolandYandex/metrics/internal/http-server/v1/metrics/getallmetrics/metrics.tmpl")
		if err != nil {
			fmt.Println("=>>>>>>>>>", err)
			log.Error("Error parsing Go template")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Prepare data for template
		var data []interface{}
		for _, metric := range metrics {
			valueStr, err := metric.ConvertValueToString()
			if err != nil {
				log.Error("Error converting metric value to string")
				continue // Skip this metric if conversion fails
			}

			data = append(data, map[string]interface{}{
				"Type":  metric.Type,
				"Name":  metric.Name,
				"Value": valueStr,
			})
		}

		// Execute the template
		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := tmpl.Execute(w, data); err != nil {
			fmt.Println("111111111111111", err)
			log.Error("Error executing Go template")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
