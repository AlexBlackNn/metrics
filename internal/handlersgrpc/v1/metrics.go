package v1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	metricsgrpc_v1 "github.com/AlexBlackNn/metrics/metricsgrpc"
	"google.golang.org/grpc"
)

type metricsService interface {
	// UpdateMetric updates a metric.
	UpdateMetric(context.Context, *metricsgrpc_v1.MetricRequest) (
		*metricsgrpc_v1.Response,
		error,
	)
	// UpdateSeveralMetrics updates several metrics.
	UpdateSeveralMetrics(context.Context, *metricsgrpc_v1.MetricsRequest) (
		*metricsgrpc_v1.Response,
		error,
	)
	// GetOneMetric gets the value of a metric.
	GetOneMetric(context.Context, *metricsgrpc_v1.MetricRequest) (
		*metricsgrpc_v1.MetricResponse,
		error,
	)
	mustEmbedUnimplementedMetricsServiceServer()
}

// serverAPI TRANSPORT layer
type serverAPI struct {
	// provides ability to work even without service interface realisation
	metricsgrpc_v1.UnimplementedMetricsServiceServer
	// service layer
	metric *metricsservice.MetricService
}

func Register(gRPC *grpc.Server, metric *metricsservice.MetricService) {
	metricsgrpc_v1.RegisterMetricsServiceServer(gRPC, &serverAPI{metric: metric})
}

func (s *serverAPI) UpdateMetric(ctx context.Context, metricgrpc *metricsgrpc_v1.MetricRequest) (
	*metricsgrpc_v1.Response,
	error,
) {
	var metricRecieved models.MetricInteraction
	var err error

	if metricgrpc.GetType() == configserver.MetricTypeCounter {
		metricRecieved, err = models.New(
			metricgrpc.GetType(),
			metricgrpc.GetId(),
			fmt.Sprintf("%d", metricgrpc.GetDelta()),
		)
	} else {
		metricRecieved, err = models.New(
			metricgrpc.GetType(),
			metricgrpc.GetId(),
			fmt.Sprintf("%g", metricgrpc.GetValue()),
		)
	}
	if err != nil {
		return &metricsgrpc_v1.Response{
			Error: "data validation failed",
		}, errors.New("data validation failed")

	}

	ctx, cancel := context.WithTimeoutCause(ctx, 300*time.Millisecond, errors.New("updateMetric timeout"))
	defer cancel()

	err = s.metric.UpdateMetricValue(ctx, metricRecieved)
	if err != nil {
		if errors.Is(err, metricsservice.ErrNotValidURL) {
			return &metricsgrpc_v1.Response{
				Error: "metric not found",
			}, errors.New("not found")
		}
		return &metricsgrpc_v1.Response{
			Error: "internal server error",
		}, errors.New("internal server error")
	}

	return &metricsgrpc_v1.Response{
		Status: "ok",
	}, nil
}

func (s *serverAPI) UpdateSeveralMetrics(context.Context, *metricsgrpc_v1.MetricsRequest) (
	*metricsgrpc_v1.Response,
	error,
) {
	return &metricsgrpc_v1.Response{
		Status: "ok",
	}, nil
}

func (s *serverAPI) GetOneMetric(context.Context, *metricsgrpc_v1.MetricRequest) (
	*metricsgrpc_v1.MetricResponse,
	error,
) {
	return &metricsgrpc_v1.MetricResponse{
		Id:    "name_counter",
		Type:  "counter",
		Delta: 3,
	}, nil
}
