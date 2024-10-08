package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AlexBlackNn/metrics/app/server"
)

var buildVersion string
var buildDate string
var buildCommit string

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

	showProjectInfo(application.Log)
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

func showProjectInfo(log *slog.Logger) {
	var sb strings.Builder

	// Определённый порядок вывода
	keys := []string{"Build version: ", "Build date: ", "Build commit: "}
	values := map[string]*string{
		"Build version: ": &buildVersion,
		"Build date: ":    &buildDate,
		"Build commit: ":  &buildCommit,
	}

	for _, key := range keys {
		if *values[key] == "" {
			*values[key] = "N/A"
		}
		sb.WriteString(key)
		sb.WriteString(*values[key])
		sb.WriteString(", ")
	}
	log.Info(strings.Trim(sb.String(), ","))
}
