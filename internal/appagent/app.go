package appagent

import (
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/transport/agenthttp"
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

	metricsService := agenthttp.New(
		log,
		cfg,
	)
	return &AppHTTP{MetricsService: metricsService}
}
