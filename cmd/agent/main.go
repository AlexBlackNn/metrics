package main

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/app_agent"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/utils"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
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
	log := utils.SetupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	monitor_application := app_agent.New(log, cfg)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitor_application.MetricsService.Start()
	}()

	wg.Add(1)
	go func() {
		for {
			time.Sleep(time.Duration(3) * time.Second)
			defer wg.Done()
			metrics := monitor_application.MetricsService.GetMetrics()
			for _, savedMetric := range metrics {
				//https://pkg.go.dev/github.com/cenkalti/backoff/v4#section-readme

				url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", savedMetric.Type, savedMetric.Name, "10")
				fmt.Println(url)
				response, err := http.Post(url, "text/pain", nil)
				if err != nil {
					fmt.Println("here")
					panic(err)
				}
				fmt.Println("==========>", response.StatusCode)
				response.Body.Close()
			}
		}
	}()

	wg.Wait()
}
