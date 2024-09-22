package v2

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/AlexBlackNn/metrics/pkg/storage/memstorage"
)

// DummyResponseWriter implements http.ResponseWriter but discards the output
type DummyResponseWriter struct {
	header http.Header
	code   int
	result []byte
	wrote  bool
}

func (w *DummyResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *DummyResponseWriter) Write(b []byte) (int, error) {
	if !w.wrote {
		w.wrote = true
	}
	copy(w.result, bytes.ReplaceAll(b, []byte("\n"), nil))
	return len(b), nil // Discard the data
}

func (w *DummyResponseWriter) WriteHeader(code int) {
	if !w.wrote {
		w.wrote = true
	}
	w.code = code
}

func ExampleReadinessProbe() {

	cfg, err := configserver.New()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Env)

	ms, err := memstorage.New(cfg, log)
	if err != nil {
		panic(err)
	}

	metricsService := metricsservice.New(
		log,
		cfg,
		ms,
		ms,
	)

	newHealthHandlers := NewHealth(log, metricsService)
	// Create a request for the benchmark
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		panic(err.Error())
	}
	w := &DummyResponseWriter{}
	// Set the header before calling ServeHTTP
	w.Header().Set("Content-Type", "json")
	w.WriteHeader(200)

	newHealthHandlers.LivenessProbe(w, req)
	fmt.Println(w.result)
	// Output:
	// Gopher
}
