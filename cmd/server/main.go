package main

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/cmd/appserver"
	"github.com/AlexBlackNn/metrics/cmd/router"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/http-server/metrics/v1"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Env)
	log.Info("starting application", slog.String("cfg", cfg.String()))

	application := appserver.New(log, cfg)
	metricsHandlers := v1.New(log, application)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	srv := &http.Server{
		Addr:         fmt.Sprintf(cfg.ServerAddr),
		Handler:      router.NewChiRouter(log, metricsHandlers),
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
