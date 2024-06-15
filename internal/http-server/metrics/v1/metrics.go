package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/cmd/appserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Metrics struct {
	log         *slog.Logger
	application *appserver.App
}

func New(log *slog.Logger, application *appserver.App) Metrics {
	return Metrics{log: log, application: application}
}

func (m *Metrics) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// TODO: in some tests somehow ClientTimeout gets 0, which creates DEADLINE ERROR
	if m.application.Cfg.ClientTimeout == 0 {
		m.application.Cfg.ClientTimeout = 10
	}
	timeout := time.Duration(m.application.Cfg.ClientTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	metrics, err := m.application.MetricsService.GetAllMetrics(ctx)

	if errors.Is(err, metricsservice.ErrMetricNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	path, err := os.Getwd()
	if err != nil {
		m.log.Error("Error getting current work dir", "err", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	pathToTemplate := filepath.Join(filepath.Dir(filepath.Dir(path)), "internal/http-server/metrics/v1/metrics.tmpl")

	tmpl, err := template.New("metrics").ParseFiles(pathToTemplate)
	if err != nil {
		m.log.Error("ParseFiles Error:", "err", err.Error(), "path:", pathToTemplate)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	var data []interface{}
	for _, metric := range metrics {
		valueStr, err := metric.ConvertValueToString()
		if err != nil {
			m.log.Error("Error converting metric value to string")
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
		m.log.Error("Error executing Go template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (m *Metrics) GetOneMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := models.CheckModelType(chi.URLParam(r, "metric_type"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: in some tests somehow ClientTimeout gets 0, which creates DEADLINE ERROR
	if m.application.Cfg.ClientTimeout == 0 {
		m.application.Cfg.ClientTimeout = 10
	}
	timeout := time.Duration(m.application.Cfg.ClientTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	metric, err := m.application.MetricsService.GetOneMetricValue(
		ctx, strings.ToLower(chi.URLParam(r, "metric_name")),
	)

	if errors.Is(err, metricsservice.ErrMetricNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v", metric.Value)))
}

func (m *Metrics) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metric, err := models.Load(
		chi.URLParam(r, "metric_type"),
		chi.URLParam(r, "metric_name"),
		chi.URLParam(r, "metric_value"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: in some tests somehow ClientTimeout gets 0, which creates DEADLINE ERROR
	if m.application.Cfg.ClientTimeout == 0 {
		m.application.Cfg.ClientTimeout = 10
	}
	timeout := time.Duration(m.application.Cfg.ClientTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = m.application.MetricsService.UpdateMetricValue(ctx, metric)
	if errors.Is(err, metricsservice.ErrNotValidURL) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
