package main

import (
	"github.com/AlexBlackNn/metrics/internal/app"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/http-server/handlers/metrics/update"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// init config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	application := app.New(log, cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	router := http.NewServeMux()
	router.HandleFunc(`/update/`, update.New(log, application))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  time.Duration(10 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(10 * time.Second),
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
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
