package v2

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/services/metricsservice"
	"github.com/AlexBlackNn/metrics/pkg/storage/memstorage"
)

// DummyResponseWriter implements http.ResponseWriter but discards the output
type DummyMetricResponseWriter struct {
	header http.Header
	code   int
	result []byte
	wrote  bool
}

func (w *DummyMetricResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *DummyMetricResponseWriter) Write(b []byte) (int, error) {
	if !w.wrote {
		w.wrote = true
	}
	w.result = make([]byte, len(b))
	n := copy(w.result, b)
	return n, nil // Discard the data
}

func (w *DummyMetricResponseWriter) WriteHeader(code int) {
	if !w.wrote {
		w.wrote = true
	}
	w.code = code
}

func ExampleMetricHandlers_UpdateMetric() {
	cfg := &configserver.Config{
		Env:                   "local",
		ServerAddr:            ":8080",
		ServerReadTimeout:     100,
		ServerWriteTimeout:    100,
		ServerIdleTimeout:     100,
		ServerRequestTimeout:  100,
		ServerStoreInterval:   100,
		ServerFileStoragePath: "/tmp/metrics-db.json",
		ServerRestore:         true,
		ServerRateLimit:       10000,
		ServerDataBaseDSN:     "DATABASE_DSN",
		ServerMigrationTable:  "migrations",
		HashKey:               "KEY",
	}

	// switch off loger output
	log := slog.New(
		slog.NewTextHandler(
			io.Discard, &slog.HandlerOptions{
				Level:     slog.LevelError,
				AddSource: true,
			},
		),
	)

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

	newMetricHandlers := New(log, metricsService)
	// Create a request for the benchmark
	body := fmt.Sprintf(`{"id":"counter_test", "type":"counter", "delta": 100}`)
	myReader := strings.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, "/value", myReader)
	if err != nil {
		panic(err.Error())
	}
	w := &DummyResponseWriter{}
	// Set the header before calling ServeHTTP
	w.Header().Set("Content-Type", "json")
	w.WriteHeader(200)

	newMetricHandlers.UpdateMetric(w, req)
	fmt.Println(string(w.result))
	// Output:
	// {"id":"counter_test","type":"counter","delta":100}
}

func ExampleMetricHandlers_GetOneMetric() {
	cfg := &configserver.Config{
		Env:                   "local",
		ServerAddr:            ":8080",
		ServerReadTimeout:     100,
		ServerWriteTimeout:    100,
		ServerIdleTimeout:     100,
		ServerRequestTimeout:  100,
		ServerStoreInterval:   100,
		ServerFileStoragePath: "/tmp/metrics-db.json",
		ServerRestore:         true,
		ServerRateLimit:       10000,
		ServerDataBaseDSN:     "DATABASE_DSN",
		ServerMigrationTable:  "migrations",
		HashKey:               "KEY",
	}

	// switch off loger output
	log := slog.New(
		slog.NewTextHandler(
			io.Discard, &slog.HandlerOptions{
				Level:     slog.LevelError,
				AddSource: true,
			},
		),
	)

	ms, err := memstorage.New(cfg, log)
	if err != nil {
		panic(err)
	}

	err = ms.UpdateMetric(context.Background(), &models.Metric[uint64]{
		Type:  "counter",
		Name:  "counter_test",
		Value: 100,
	})
	if err != nil {
		panic(err)
	}

	metricsService := metricsservice.New(
		log,
		cfg,
		ms,
		ms,
	)

	newMetricHandlers := New(log, metricsService)
	// Create a request for the benchmark
	body := fmt.Sprintf(`{"id":"counter_test", "type":"counter", "delta": 100}`)
	myReader := strings.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, "/value", myReader)
	if err != nil {
		panic(err.Error())
	}
	w := &DummyResponseWriter{}
	// Set the header before calling ServeHTTP
	w.Header().Set("Content-Type", "json")
	w.WriteHeader(200)

	newMetricHandlers.GetOneMetric(w, req)
	fmt.Println(string(w.result))
	// Output:
	// {"id":"counter_test","type":"counter","delta":100}
}
