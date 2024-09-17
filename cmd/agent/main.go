package main

import (
	"context"
	"github.com/AlexBlackNn/metrics/app/agent"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup

	cfg, err := configagent.New()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	appMonitor := agent.NewAppMonitor(log, cfg)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	wg.Add(4)
	go func() {
		defer wg.Done()
		<-stop
		cancel()
	}()

	go func() {
		defer wg.Done()
		appMonitor.MetricsService.Collect(ctx)
	}()

	go func() {
		defer wg.Done()
		appMonitor.MetricsService.CollectAddition(ctx)
	}()

	go func() {
		defer wg.Done()
		appMonitor.MetricsService.Send(ctx)
	}()

	wg.Wait()
}
