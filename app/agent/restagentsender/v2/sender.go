package v2

import (
	"context"
	"fmt"
	"github.com/AlexBlackNn/metrics/app/agent/hash"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/agentmetricsservice"
	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
	"log/slog"
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
	reportInterval := time.Duration(s.cfg.ReportInterval) * time.Second
	rateLimiter := rate.NewLimiter(rate.Limit(s.cfg.AgentRateLimit), s.cfg.AgentBurstTokens)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(reportInterval):
			for _, savedMetric := range s.GetMetrics() {
				err := rateLimiter.Wait(ctx)
				if err != nil {
					log.Error(err.Error())
					return
				}
				go func(savedMetric models.MetricInteraction) {
					restyClient := resty.New()
					restyClient.
						SetRetryCount(s.cfg.AgentRetryCount).
						SetRetryWaitTime(s.cfg.AgentRetryWaitTime).
						SetRetryMaxWaitTime(s.cfg.AgentRetryMaxWaitTime)

					var body string
					if savedMetric.GetType() == configserver.MetricTypeCounter {
						body = fmt.Sprintf(`{"id":"%s", "type":"%s", "delta": %d}`,
							savedMetric.GetName(),
							savedMetric.GetType(),
							savedMetric.GetValue(),
						)
					} else {
						body = fmt.Sprintf(`{"id":"%s", "type":"%s", "value": %v}`,
							savedMetric.GetName(),
							savedMetric.GetType(),
							savedMetric.GetValue(),
						)
					}

					hashCalculator := hash.New(s.cfg)
					hashResult := hashCalculator.MetricHash(body)

					url := fmt.Sprintf("http://%s/update/", s.cfg.ServerAddr)
					log.Info("sending data", "url", url)

					resp, err := restyClient.R().
						SetHeader("Content-Type", "application/json").
						SetHeader("HashSHA256", hashResult).
						SetBody(body).
						Post(url)
					if err != nil {
						log.Error("error creating http request")
						return
					}
					log.Info("http request finished successfully",
						"url", url,
						"statusCode", resp.StatusCode(),
						"body", string(resp.Body()),
					)
				}(savedMetric)
			}
		}
	}
}
