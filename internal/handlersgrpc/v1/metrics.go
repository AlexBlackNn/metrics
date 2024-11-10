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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Errorf(codes.InvalidArgument, "data validation failed")

	}

	ctx, cancel := context.WithTimeoutCause(ctx, 300*time.Millisecond, errors.New("updateMetric timeout"))
	defer cancel()

	err = s.metric.UpdateMetricValue(ctx, metricRecieved)
	if err != nil {
		if errors.Is(err, metricsservice.ErrNotValidURL) {
			return nil, status.Errorf(codes.NotFound, "data not found")
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &metricsgrpc_v1.Response{
		Status: "ok",
	}, nil
}

func (s *serverAPI) UpdateSeveralMetrics(ctx context.Context, severalMetricsgrpc *metricsgrpc_v1.MetricsRequest) (
	*metricsgrpc_v1.Response,
	error,
) {
	var metricRecieved models.MetricInteraction
	var err error

	severalMetrics := make([]models.MetricInteraction, len(severalMetricsgrpc.GetMetrics()))

	for i, metricgrpc := range severalMetricsgrpc.GetMetrics() {
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
			return nil, status.Errorf(codes.InvalidArgument, "data validation failed")
		}
		severalMetrics[i] = metricRecieved
	}
	ctx, cancel := context.WithTimeoutCause(ctx, 300*time.Millisecond, errors.New("updateMetric timeout"))
	defer cancel()

	err = s.metric.UpdateSeveralMetrics(ctx, severalMetrics)
	if err != nil {
		if errors.Is(err, metricsservice.ErrNotValidURL) {
			return nil, status.Errorf(codes.NotFound, "data not found")
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &metricsgrpc_v1.Response{
		Status: "ok",
	}, nil
}

func (s *serverAPI) GetOneMetric(ctx context.Context, metricgrpc *metricsgrpc_v1.MetricRequest) (
	*metricsgrpc_v1.MetricResponse,
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
		return nil, status.Errorf(codes.InvalidArgument, "data validation failed")
	}

	ctx, cancel := context.WithTimeoutCause(
		ctx,
		300*time.Millisecond,
		errors.New("updateMetric timeout"),
	)
	defer cancel()

	metricReturned, err := s.metric.GetOneMetricValue(ctx, metricRecieved)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "data extraction failed")
	}
	if metricReturned.GetType() == configserver.MetricTypeCounter {
		delta, ok := metricReturned.GetValue().(uint64)
		if !ok {
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
		return &metricsgrpc_v1.MetricResponse{
			Id:    metricReturned.GetName(),
			Type:  metricReturned.GetType(),
			Delta: int64(delta),
		}, nil
	}
	value, ok := metricReturned.GetValue().(float64)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &metricsgrpc_v1.MetricResponse{
		Id:    metricReturned.GetName(),
		Type:  metricReturned.GetType(),
		Value: value,
	}, nil

}
