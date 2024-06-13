package getallmetrics

import (
	"context"
	"errors"
	"github.com/AlexBlackNn/metrics/internal/appserver"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func New(log *slog.Logger, application *appserver.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// TODO: in some tests somehow ClientTimeout gets 0, which creates DEADLINE ERROR
		if application.Cfg.ClientTimeout == 0 {
			application.Cfg.ClientTimeout = 10
		}
		timeout := time.Duration(application.Cfg.ClientTimeout) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		metrics, err := application.MetricsService.GetAllMetrics(ctx)

		if errors.Is(err, metricsservice.ErrMetricNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		path, err := os.Getwd()
		if err != nil {
			log.Error("Error getting current work dir", "err", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		pathToTemplate := filepath.Join(filepath.Dir(filepath.Dir(path)), "internal/http-server/v1/metrics/getallmetrics/metrics.tmpl")

		tmpl, err := template.New("metrics").ParseFiles(pathToTemplate)
		if err != nil {
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

		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := tmpl.Execute(w, data); err != nil {
			log.Error("Error executing Go template")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
