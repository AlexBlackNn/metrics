package main

import (
	"context"
	"github.com/AlexBlackNn/metrics/app/agent"
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

	appHTTP := agent.NewAppHTTP(log, cfg)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	wg.Add(3)
	go func() {
		defer wg.Done()
		<-stop
		cancel()
	}()

	go func() {
		defer wg.Done()
		appHTTP.MetricsService.Start(ctx)
	}()

	go func() {
		defer wg.Done()
		appHTTP.MetricsService.Send(ctx)
	}()
	wg.Wait()
}
