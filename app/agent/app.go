package agent

import (
	"context"
	"github.com/AlexBlackNn/metrics/app/agent/restagentsender"
	"github.com/AlexBlackNn/metrics/internal/config"
	"log/slog"
)

type CollectSender interface {
	Collect(ctx context.Context)
	Send(ctx context.Context)
}

// AppMonitor service consists all service layers
type AppMonitor struct {
	MetricsService CollectSender
}

// NewAppMonitor creates App
func NewAppMonitor(
	log *slog.Logger,
	cfg *config.Config,
) *AppMonitor {

	metricsService := restagentsender.New(
		log,
		cfg,
	)
	return &AppMonitor{MetricsService: metricsService}
}
