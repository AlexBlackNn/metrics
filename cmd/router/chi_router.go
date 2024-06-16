package router

import (
	"github.com/AlexBlackNn/metrics/internal/api/metrics/v1"
	projectLogger "github.com/AlexBlackNn/metrics/internal/api/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func NewChiRouter(log *slog.Logger, m v1.Metrics) chi.Router {

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(projectLogger.New(log))
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/", m.GetAllMetrics)
		r.Post("/update/{metric_type}/{metric_name}/{metric_value}", m.UpdateMetric)
		r.Get("/value/{metric_type}/{metric_name}", m.GetOneMetric)
	})
	return router
}
