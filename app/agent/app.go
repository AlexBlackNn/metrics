package agent

import (
	"context"
	"github.com/AlexBlackNn/metrics/app/agent/restagentsender/v2"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"log/slog"
)

type CollectSender interface {
	Collect(ctx context.Context)
	Send(ctx context.Context)
}

// AppMonitor service consists all service layers.
type AppMonitor struct {
	MetricsService CollectSender
}

// NewAppMonitor creates App.
func NewAppMonitor(
	log *slog.Logger,
	cfg *configagent.Config,
) *AppMonitor {

	metricsService := v2.New(
		log,
		cfg,
	)
	return &AppMonitor{MetricsService: metricsService}
}
