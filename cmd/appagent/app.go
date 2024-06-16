package appagent

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/sender/restagentsender"
	"log/slog"
)

type AgentService interface {
	Start(<-chan struct{})
	Send(<-chan struct{})
}

// AppHTTP service consists all service layers
type AppHTTP struct {
	MetricsService AgentService
}

// NewAppHTTP creates App
func NewAppHTTP(
	log *slog.Logger,
	cfg *config.Config,
) *AppHTTP {

	metricsService := restagentsender.New(
		log,
		cfg,
	)
	return &AppHTTP{MetricsService: metricsService}
}
