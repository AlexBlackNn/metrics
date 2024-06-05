package main

import (
	"github.com/AlexBlackNn/metrics/internal/app_agent"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/utils"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// init config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	// init logger
	log := utils.SetupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	appHttp := app_agent.NewAppHttp(log, cfg)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		appHttp.MetricsService.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		appHttp.MetricsService.Transmit()
	}()
	wg.Wait()
}
