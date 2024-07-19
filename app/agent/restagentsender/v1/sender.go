package v1

import (
	"context"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/agentmetricsservice"
	"log/slog"
	"net/http"
	"time"
)

type Sender struct {
	log *slog.Logger
	cfg *configagent.Config
	*agentmetricsservice.MonitorService
}

func New(
	log *slog.Logger,
	cfg *configagent.Config,
) *Sender {
	return &Sender{
		log,
		cfg,
		agentmetricsservice.New(log, cfg),
	}
}

func (s *Sender) Send(ctx context.Context) {

	log := s.log.With(
		slog.String("info", "SERVICE LAYER: metricsHttpService.Transmit"),
	)
	client := http.Client{
		Timeout: time.Duration(s.cfg.AgentTimeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { // в 1 инкрименте "Редиректы не поддерживаются."
			return http.ErrUseLastResponse
		}}
	reportInterval := time.Duration(s.cfg.ReportInterval) * time.Second
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, savedMetric := range s.GetMetrics() {
				go func(savedMetric models.MetricInteraction) {
					savedMetricValue := savedMetric.GetStringValue()

					url := fmt.Sprintf("http://%s/update/%s/%s/%s", s.cfg.ServerAddr, savedMetric.GetType(), savedMetric.GetName(), savedMetricValue)
					log.Info("sending data", "url", url)

					req, err := http.NewRequest(http.MethodPost, url, nil) // (1)

					// TODO: need refactoring to better work with error.
					if err != nil {
						log.Error("error creating http request")
						return
					}

					//Would be better to add backoff, but in next task client itself can do it.
					//https://pkg.go.dev/github.com/cenkalti/backoff/v4#section-readme
					response, err := client.Do(req)

					// TODO: need refactoring to better work with error.
					if err != nil {
						log.Error("error doing http request", "err", err.Error())
						return
					}
					log.Info("sending data", "url", url, "status_code", response.StatusCode)
					response.Body.Close()

				}(savedMetric)
			}
			<-time.After(reportInterval)
		}
	}
}
