package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexBlackNn/metrics/app/server"
)

// @title           Swagger API
// @version         1.0
// @description     metric collection service.
// @contact.name   API Support
// @license.name  Apache 2.0
// @license.calculation   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
//
//go:generate go run github.com/swaggo/swag/cmd/swag init
func main() {

	application, err := server.New()
	if err != nil {
		panic(err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	application.Log.Info("starting application", slog.String("cfg", application.Cfg.String()))
	go func() {
		if err = application.Srv.ListenAndServe(); err != nil {
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
