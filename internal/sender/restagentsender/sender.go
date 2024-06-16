package restagentsender

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/agentmetricsservice"
	"log/slog"
	"net/http"
	"time"
)

type Sender struct {
	log *slog.Logger
	cfg *config.Config
	*agentmetricsservice.MetricsService
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *Sender {
	return &Sender{
		log,
		cfg,
		agentmetricsservice.New(log, cfg),
	}
}

func (mhs *Sender) Send(stop <-chan struct{}) {

	log := mhs.log.With(
		slog.String("info", "SERVICE LAYER: metricsHttpService.Transmit"),
	)
	client := http.Client{
		Timeout: time.Duration(mhs.cfg.ClientTimeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { // в 1 инкрименте "Редиректы не поддерживаются."
			return http.ErrUseLastResponse
		}}
	for {
		select {
		case <-stop:
			return
		default:
			for _, savedMetric := range mhs.GetMetrics() {
				go func(savedMetric models.MetricInteraction) {
					savedMetricValue := savedMetric.GetStringValue()

					url := fmt.Sprintf("http://%s/update/%s/%s/%s", mhs.cfg.ServerAddr, savedMetric.GetType(), savedMetric.GetName(), savedMetricValue)
					log.Info("sending data", "url", url)

					req, err := http.NewRequest(http.MethodPost, url, nil) // (1)

					// TODO: need refactoring to better work with error
					if err != nil {
						log.Error("error creating http request")
					}

					//// TODO: find out why without it EOF?
					////https://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi
					//req.Close = true

					// would be better to add backoff, but in next task client itself can do it
					//https://pkg.go.dev/github.com/cenkalti/backoff/v4#section-readme
					response, err := client.Do(req)

					// TODO: need refactoring to better work with error
					if err != nil {
						log.Error("error doing http request", "err", err.Error())
					} else {
						log.Info("sending data", "url", url, "status_code", response.StatusCode)
						response.Body.Close()
					}
				}(savedMetric)
			}
			<-time.After(time.Duration(mhs.cfg.ReportInterval) * time.Second)
		}
	}
}
