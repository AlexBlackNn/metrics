package v1

import (
	"context"

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
