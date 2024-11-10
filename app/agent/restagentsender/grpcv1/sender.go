package grpcv1

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AlexBlackNn/metrics/app/agent/encryption"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/agentmetricsservice"
	metricsgrpc_v1 "github.com/AlexBlackNn/metrics/metricsgrpc"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

type Sender struct {
	log       *slog.Logger
	cfg       *configagent.Config
	encryptor *encryption.Encryptor
	*agentmetricsservice.MonitorService
	conn   *grpc.ClientConn
	client metricsgrpc_v1.MetricsServiceClient
}

func New(
	log *slog.Logger,
	cfg *configagent.Config,
	encryptor *encryption.Encryptor,
) *Sender {

	conn, err := grpc.NewClient(cfg.ServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	client := metricsgrpc_v1.NewMetricsServiceClient(conn)

	return &Sender{
		log,
		cfg,
		encryptor,
		agentmetricsservice.New(log, cfg),
		conn,
		client,
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
					req := &metricsgrpc_v1.MetricRequest{
						Id:   savedMetric.GetName(),
						Type: savedMetric.GetType(),
					}
					if savedMetric.GetType() == configserver.MetricTypeCounter {
						switch v := savedMetric.GetValue().(type) {
						case int64:
							req.Delta = v
						case uint64:
							req.Delta = int64(v)
						}
					} else {
						switch v := savedMetric.GetValue().(type) {
						case int64:
							req.Value = float64(v)
						case float64:
							req.Value = v
						case uint64:
							req.Value = float64(v)
						case int32:
							req.Value = float64(v)
						case uint32:
							req.Value = float64(v)
						default:
							log.Error("unsupported metric type", "type", fmt.Sprintf("%T", v))
							return
						}
					}

					_, err := s.client.UpdateMetric(ctx, req)
					if err != nil {
						log.Error("error sending metric", "error", err.Error())
						return
					}
					log.Info("metric sent successfully")
				}(savedMetric)
			}
		}
	}
}
