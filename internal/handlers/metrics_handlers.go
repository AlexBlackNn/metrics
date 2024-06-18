package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/app/server"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MetricHandlers struct {
	application *server.App
}

func New(application *server.App) MetricHandlers {
	return MetricHandlers{application: application}
}

func (m *MetricHandlers) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	metrics, err := m.application.MetricsService.GetAllMetrics(ctx)

	if errors.Is(err, metricsservice.ErrMetricNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	path, err := os.Getwd()
	if err != nil {
		m.application.Log.Error("Error getting current work dir", "err", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	pathToTemplate := filepath.Join(filepath.Dir(filepath.Dir(path)), "internal/handlers/metrics.tmpl")

	tmpl, err := template.New("metrics").ParseFiles(pathToTemplate)
	if err != nil {
		m.application.Log.Error("ParseFiles Error:", "err", err.Error(), "path:", pathToTemplate)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	var data []interface{}
	for _, metric := range metrics {
		data = append(data, map[string]interface{}{
			"Type":  metric.GetType(),
			"Name":  metric.GetName(),
			"Value": metric.GetStringValue(),
		})
	}
	w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, data); err != nil {
		m.application.Log.Error("Error executing Go template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (m *MetricHandlers) GetOneMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := models.CheckModelType(chi.URLParam(r, "metric_type"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
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
	w.Write([]byte(fmt.Sprintf("%v", metric.GetValue())))
}

func (m *MetricHandlers) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metric, err := models.New(
		chi.URLParam(r, "metric_type"),
		chi.URLParam(r, "metric_name"),
		chi.URLParam(r, "metric_value"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
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
