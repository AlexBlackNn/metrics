package main

import (
	"github.com/AlexBlackNn/metrics/cmd/appagent"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	appHTTP := appagent.NewAppHTTP(log, cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	cancel := make(chan struct{})

	wg.Add(3)
	go func() {
		defer wg.Done()
		<-stop
		close(cancel)
	}()

	go func() {
		defer wg.Done()
		appHTTP.MetricsService.Start(cancel)
	}()

	go func() {
		defer wg.Done()
		appHTTP.MetricsService.Transmit(cancel)
	}()
	wg.Wait()
}
