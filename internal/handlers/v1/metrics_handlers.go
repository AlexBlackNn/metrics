package v1

import (
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type MetricHandlers struct {
	log            *slog.Logger
	metricsService *metricsservice.MetricService
}

func New(log *slog.Logger, metricsService *metricsservice.MetricService) MetricHandlers {
	return MetricHandlers{log: log, metricsService: metricsService}
}

func (m *MetricHandlers) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	metrics, err := m.metricsService.GetAllMetrics(ctx)
	if errors.Is(err, metricsservice.ErrMetricNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	path, err := os.Getwd()
	if err != nil {
		m.log.Error("Error getting current work dir", "err", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	pathToTemplate := filepath.Join(path, "internal/handlers/v1/metrics.tmpl")
	tmpl, err := template.New("metrics").ParseFiles(pathToTemplate)
	if err != nil {
		m.log.Error("ParseFiles Error:", "err", err.Error(), "path:", pathToTemplate)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, metrics); err != nil {
		m.log.Error("Error executing Go template")
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

	ctx := r.Context()
	metric, err := models.New(chi.URLParam(r, "metric_type"), chi.URLParam(r, "metric_name"), "0")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricReturned, err := m.metricsService.GetOneMetricValue(ctx, metric)
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
	w.Write([]byte(fmt.Sprintf("%v", metricReturned.GetValue())))
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

	ctx := r.Context()
	err = m.metricsService.UpdateMetricValue(ctx, metric)
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
