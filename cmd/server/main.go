package main

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/appserver"
	"github.com/AlexBlackNn/metrics/internal/config"
	projectLogger "github.com/AlexBlackNn/metrics/internal/http-server/middleware/logger"
	"github.com/AlexBlackNn/metrics/internal/http-server/v1/metrics/getallmetrics"
	"github.com/AlexBlackNn/metrics/internal/http-server/v1/metrics/getonemetric"
	"github.com/AlexBlackNn/metrics/internal/http-server/v1/metrics/update"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewChiRouter(log *slog.Logger, application *appserver.App) chi.Router {

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(projectLogger.New(log))
	router.Use(middleware.Recoverer)
	//router.Use(middleware.URLFormat)

	//router.Route(`/update/`, update.New(log, application))

	router.Route("/update/", func(r chi.Router) {
		r.Post("/{metric_type}/{metric_name}/{metric_value}", update.New(log, application))
		//r.Get("/", expression.New(log, application))
	})
	router.Route("/value/", func(r chi.Router) {
		r.Get("/{metric_type}/{metric_name}", getonemetric.New(log, application))
		//r.Get("/", expression.New(log, application))
	})
	router.Get("/", getallmetrics.New(log, application))
	return router
}

func main() {
	// init config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	application := appserver.New(log, cfg)
	router := NewChiRouter(log, application)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	srv := &http.Server{
		Addr:         fmt.Sprintf(cfg.ServerAddr),
		Handler:      router,
		ReadTimeout:  time.Duration(10 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(10 * time.Second),
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	log.Info("server started")

	signalType := <-stop
	log.Info(
		"application stopped",
		slog.String("signalType",
			signalType.String()),
	)

}

const (
	envLocal = "local"
	envDemo  = "demo"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	case envDemo:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{
					Level:     slog.LevelInfo,
					AddSource: true,
				},
			),
		)
	}
	return log
}
