package middleware

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

type GzipWriter struct {
	ResWriter       http.ResponseWriter
	Writer          *gzip.Writer
	GzipWriterMutex sync.Mutex
}

func (w *GzipWriter) Header() http.Header {
	return w.ResWriter.Header()
}

func (w *GzipWriter) WriteHeader(statusCode int) {
	if !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "application/json") &&
		!strings.Contains(w.ResWriter.Header().Get("Content-Type"), "text/html") {
		w.ResWriter.WriteHeader(statusCode)
		return
	}
	w.ResWriter.Header().Set("Content-Encoding", "gzip")
	w.ResWriter.WriteHeader(statusCode)
}

func (w *GzipWriter) Write(b []byte) (int, error) {
	if !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "application/json") &&
		!strings.Contains(w.ResWriter.Header().Get("Content-Type"), "text/html") {
		return w.ResWriter.Write(b)
	}
	defer func(Writer *gzip.Writer) {
		err := Writer.Flush()
		if err != nil {
			io.WriteString(w, err.Error())
		}

	}(w.Writer)

	return w.Writer.Write(b)
}

func (w *GzipWriter) Close() error {
	return w.Writer.Close()
}

func GzipCompressor(log *slog.Logger, compressorLevel int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/gzip"),
		)
		log.Info("gzip compressor enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(strings.Join(r.Header.Values("Accept-Encoding"), " "), "gzip") {
				// if gzip is not supported then return uncompressed page
				next.ServeHTTP(w, r)
				return
			}
			log.Info("gzip is supported")

			gzipWr, err := gzip.NewWriterLevel(w, compressorLevel)
			defer gzipWr.Close()
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			next.ServeHTTP(&GzipWriter{ResWriter: w, Writer: gzipWr}, r)
		}
		return http.HandlerFunc(fn)
	}
}
