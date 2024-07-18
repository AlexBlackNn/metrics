package v3

import (
	"context"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/services/agentmetricsservice"
	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
	"log/slog"
	"strings"
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

			var body strings.Builder

			body.WriteString("[")
			for _, savedMetric := range s.GetMetrics() {
				if savedMetric.GetType() == configserver.MetricTypeCounter {
					body.WriteString(fmt.Sprintf(`{"id":"%s", "type":"%s", "delta": %d},`,
						savedMetric.GetName(),
						savedMetric.GetType(),
						savedMetric.GetValue(),
					))
				} else {
					body.WriteString(fmt.Sprintf(`{"id":"%s", "type":"%s", "value": %v},`,
						savedMetric.GetName(),
						savedMetric.GetType(),
						savedMetric.GetValue(),
					))
				}
			}
			// need to delete last comma
			jsonBody := body.String()[:len(body.String())-1] + "]"

			err := rateLimiter.Wait(ctx)
			if err != nil {
				log.Error(err.Error())
				return
			}

			restyClient := resty.New()
			restyClient.
				SetRetryCount(s.cfg.AgentRetryCount).
				SetRetryWaitTime(s.cfg.AgentRetryWaitTime).
				SetRetryMaxWaitTime(s.cfg.AgentRetryMaxWaitTime)

			url := fmt.Sprintf("http://%s/updates/", s.cfg.ServerAddr)
			log.Info("sending data", "url", url)
			resp, err := restyClient.R().
				SetHeader("Content-Type", "application/json").
				SetBody(jsonBody).
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
		}
	}
}
