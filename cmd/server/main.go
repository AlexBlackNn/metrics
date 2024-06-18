package main

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/app/server"
	"github.com/AlexBlackNn/metrics/cmd/server/router"
	"github.com/AlexBlackNn/metrics/internal/handlers"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	application, err := server.New()
	if err != nil {
		panic(err)
	}
	metricsHandlers := handlers.New(application)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	srv := &http.Server{
		Addr:         fmt.Sprintf(application.Cfg.ServerAddr),
		Handler:      router.NewChiRouter(application.Log, metricsHandlers),
		ReadTimeout:  time.Duration(application.Cfg.ServerReadTimeout) * time.Second,
		WriteTimeout: time.Duration(application.Cfg.ServerWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(application.Cfg.ServerIdleTimeout) * time.Second,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	application.Log.Info("server started")

	signalType := <-stop
	application.Log.Info(
		"application stopped",
		slog.String("signalType",
			signalType.String()),
	)

}
