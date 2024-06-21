package router

import (
	"github.com/AlexBlackNn/metrics/internal/handlers/v1"
	"github.com/AlexBlackNn/metrics/internal/handlers/v2"

	customMiddleware "github.com/AlexBlackNn/metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func NewChiRouter(log *slog.Logger, metricHandlerV1 v1.MetricHandlers, metricHandlerV2 v2.MetricHandlers) *chi.Mux {

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(customMiddleware.Logger(log))
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/", metricHandlerV1.GetAllMetrics)
		r.Post("/update/{metric_type}/{metric_name}/{metric_value}", metricHandlerV1.UpdateMetric)
		r.Get("/value/{metric_type}/{metric_name}", metricHandlerV1.GetOneMetric)
		r.Post("/update/", metricHandlerV2.UpdateMetric)
		r.Post("/value/", metricHandlerV2.GetOneMetric)
	})
	return router
}
