package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// content type to compress data
var compressibleContentTypes = []string{
	"application/json",
	"text/html",
}

type gzipWriter struct {
	ResWriter http.ResponseWriter
	Writer    io.Writer
}

func (w *gzipWriter) Header() http.Header {
	return w.ResWriter.Header()
}

func (w *gzipWriter) WriteHeader(statusCode int) {
	if !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "application/json") && !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "text/html") {
		fmt.Println(11111)
		w.ResWriter.WriteHeader(statusCode)
	}
	fmt.Println(222222)
	w.ResWriter.Header().Set("Content-Encoding", "gzip")
	w.ResWriter.WriteHeader(statusCode)
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	if !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "application/json") && !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "text/html") {
		fmt.Println(33333)
		return w.ResWriter.Write(b)
	}
	fmt.Println(4444)
	return w.Writer.Write(b)
}

func GzipCompressor(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/gzip"),
		)

		log.Info("gzip middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(strings.Join(r.Header.Values("Accept-Encoding"), " "), "gzip") {
				// if gzip is not supported then return uncompressed page
				next.ServeHTTP(w, r)
				return
			}
			log.Info("gzip is supported")
			gzipWr, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gzipWr.Close()

			next.ServeHTTP(&gzipWriter{ResWriter: w, Writer: gzipWr}, r)
		}

		return http.HandlerFunc(fn)
	}
}
