package agentmetricsservice

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type MetricsHttpService struct {
	log *slog.Logger
	cfg *config.Config
	*MetricsService
}

func NewMetricsHttpService(
	log *slog.Logger,
	cfg *config.Config,
) *MetricsHttpService {
	return &MetricsHttpService{
		log,
		cfg,
		New(log, cfg),
	}
}

func (mhs *MetricsHttpService) Transmit() {

	var wg sync.WaitGroup

	//TODO:// get timeout from config
	client := http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { // в 1 инкрименте "Редиректы не поддерживаются."
			return http.ErrUseLastResponse
		}}

	for {
		time.Sleep(time.Duration(3) * time.Second)

		metrics := mhs.GetMetrics()

		wg.Add(len(metrics))
		for _, savedMetric := range metrics {
			go func(savedMetric models.Metric) {
				defer wg.Done()
				// TODO: convert any (int64, float64,...) to string
				// TODO: backoff
				//https://pkg.go.dev/github.com/cenkalti/backoff/v4#section-readme
				url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", savedMetric.Type, savedMetric.Name, "10")
				fmt.Println(url)

				req, err := http.NewRequest(http.MethodPost, url, nil) // (1)
				// TODO: find out why without it EOF?
				//https://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi
				req.Close = true
				if err != nil {
					//TODO: HANDLE error do not panic
					panic(err)
				}

				response, err := client.Do(req)
				if err != nil {
					//TODO: HANDLE error do not exit
					fmt.Println("11111111111111", err)
					os.Exit(1)
				}

				fmt.Println("==========>", response.StatusCode)
				response.Body.Close()
			}(savedMetric)
		}
	}
}
