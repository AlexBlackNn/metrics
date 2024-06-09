package main

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/appserver"
	"github.com/AlexBlackNn/metrics/internal/config"
	projectLogger "github.com/AlexBlackNn/metrics/internal/http-server/middleware/logger"
	"github.com/AlexBlackNn/metrics/internal/http-server/v1/metrics/getallmetrics"
	"github.com/AlexBlackNn/metrics/internal/http-server/v1/metrics/getonemetric"
	"github.com/AlexBlackNn/metrics/internal/http-server/v1/metrics/updatemetric"
	"github.com/AlexBlackNn/metrics/internal/utils"
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

	return router.Route("/", func(r chi.Router) {
		r.Get("/", getallmetrics.New(log, application))
		r.Post("/update/{metric_type}/{metric_name}/{metric_value}", updatemetric.New(log, application))
		r.Get("/value/{metric_type}/{metric_name}", getonemetric.New(log, application))
	})
}

func main() {
	// init config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	// init logger
	log := utils.SetupLogger(cfg.Env)
	log.Info("starting application", slog.String("cfg", cfg.String()))
	application := appserver.New(log, cfg)
	router := NewChiRouter(log, application)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	srv := &http.Server{
		Addr:         fmt.Sprintf(cfg.ServerAddr),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ServerReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.ServerIdleTimeout) * time.Second,
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
