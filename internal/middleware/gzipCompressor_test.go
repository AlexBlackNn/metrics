package middleware

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"testing"
)

// DummyResponseWriter implements http.ResponseWriter but discards the output
type DummyResponseWriter struct {
	header http.Header
	code   int
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
	return len(b), nil // Discard the data
}

func (w *DummyResponseWriter) WriteHeader(code int) {
	if !w.wrote {
		w.wrote = true
	}
	w.code = code
}

func BenchmarkGzipCompressor(b *testing.B) {
	// Create a dummy handler to simulate a real handler
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate some data to be written
		_, err := io.WriteString(w, "This is some test data to be compressed.")
		if err != nil {
			b.Fatal(err)
		}
	})

	// Create a logger for the benchmark
	logger := slog.New(
		slog.NewTextHandler(
			io.Discard, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
			},
		),
	)

	// Create a request for the benchmark
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	// Run the benchmark
	gzipCompressor := GzipCompressor(logger, gzip.DefaultCompression)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a DummyResponseWriter to discard output
		w := &DummyResponseWriter{}

		// Set the header before calling ServeHTTP
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(200)

		gzipCompressor(dummyHandler).ServeHTTP(w, req)
	}
}
