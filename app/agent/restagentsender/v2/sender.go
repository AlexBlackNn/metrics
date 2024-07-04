package v2

import (
	"context"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/agentmetricsservice"
	"github.com/go-resty/resty/v2"
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

func (mhs *Sender) Send(ctx context.Context) {

	log := mhs.log.With(
		slog.String("info", "SERVICE LAYER: metricsHttpService.Transmit"),
	)
	reportInterval := time.Duration(mhs.cfg.ReportInterval) * time.Second
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, savedMetric := range mhs.GetMetrics() {
				go func(savedMetric models.MetricInteraction) {
					restyClient := resty.New()
					restyClient.
						SetRetryCount(mhs.cfg.AgentRetryCount).
						SetRetryWaitTime(mhs.cfg.AgentRetryWaitTime).
						SetRetryMaxWaitTime(mhs.cfg.AgentRetryMaxWaitTime)

					var body string
					if savedMetric.GetType() == configserver.Counter {
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
					url := fmt.Sprintf("http://%s/update/", mhs.cfg.ServerAddr)
					log.Info("sending data", "url", url)
					resp, err := restyClient.R().
						SetHeader("Content-Type", "application/json").
						SetBody(body).
						Post(url)
					if err != nil {
						log.Error("error creating http request")
					}
					log.Info("http request finished successfully",
						"url", url,
						"statusCode", resp.StatusCode(),
						"body", string(resp.Body()),
					)
				}(savedMetric)
			}
			<-time.After(reportInterval)
		}
	}
}
