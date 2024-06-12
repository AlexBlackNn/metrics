package appagent

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/services/agentmetricsservice"
	"log/slog"
)

type AgentServiceInterface interface {
	Start(chan struct{})
	Transmit(chan struct{})
}

// AppHTTP service consists all service layers
type AppHTTP struct {
	MetricsService AgentServiceInterface
}

// NewAppHTTP creates App
func NewAppHTTP(
	log *slog.Logger,
	cfg *config.Config,
) *AppHTTP {

	metricsService := agentmetricsservice.NewHTTPService(
		log,
		cfg,
	)
	return &AppHTTP{MetricsService: metricsService}
}
