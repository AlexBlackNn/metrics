package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/AlexBlackNn/metrics/app/agent"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"github.com/AlexBlackNn/metrics/internal/logger"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	var wg sync.WaitGroup

	cfg, err := configagent.New()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Env)
	showProjectInfo(log)
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

func showProjectInfo(log *slog.Logger) {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	projInfo := fmt.Sprintf(
		"Build version: %s, Build date: %s, Build commit: %s",
		buildVersion, buildDate, buildCommit,
	)
	log.Info(projInfo)
}
