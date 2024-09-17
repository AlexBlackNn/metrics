package router

import (
	"compress/gzip"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/handlers/v1"
	"github.com/AlexBlackNn/metrics/internal/handlers/v2"
	v3 "github.com/AlexBlackNn/metrics/internal/handlers/v3"
	customMiddleware "github.com/AlexBlackNn/metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"log/slog"
	"time"
)

func NewChiRouter(
	cfg *configserver.Config,
	log *slog.Logger,
	metricHandlerV1 v1.MetricHandlers,
	metricHandlerV2 v2.MetricHandlers,
	healthHandlerV2 v2.HealthHandlers,
	metricHandlerV3 v3.MetricHandlers,
) *chi.Mux {

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	//	Rate limit by IP and URL path (aka endpoint)
	router.Use(httprate.Limit(
		cfg.ServerRateLimit, // requests
		time.Second,         // per duration
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))
	router.Use(customMiddleware.Logger(log))
	router.Use(customMiddleware.HashChecker(log, cfg))
	router.Use(customMiddleware.GzipDecompressor(log))
	router.Use(customMiddleware.GzipCompressor(log, gzip.BestCompression))

	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/", metricHandlerV1.GetAllMetrics)
		r.Get("/ping", healthHandlerV2.ReadinessProbe)
		r.Post("/update/{metric_type}/{metric_name}/{metric_value}", metricHandlerV1.UpdateMetric)
		r.Get("/value/{metric_type}/{metric_name}", metricHandlerV1.GetOneMetric)
		r.Post("/update/", metricHandlerV2.UpdateMetric)
		r.Post("/updates/", metricHandlerV3.UpdateSeveralMetrics)
		r.Post("/value/", metricHandlerV2.GetOneMetric)
	})
	return router
}
