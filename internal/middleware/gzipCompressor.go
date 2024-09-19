package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

type gzipWriterPool struct {
	p sync.Pool
}

var gzipWrPool *gzipWriterPool

func init() {
	gzipWrPool = &gzipWriterPool{}
}

func (gp *gzipWriterPool) Get(w http.ResponseWriter, compressorLevel int) (*GzipWriter, error) {

	gzipWriter := gp.p.Get()
	if gzipWriter == nil {
		gzipWr, err := gzip.NewWriterLevel(w, compressorLevel)
		if err != nil {
			return nil, err
		}
		return &GzipWriter{ResWriter: w, Writer: gzipWr}, nil
	}
	return gzipWriter.(*GzipWriter), nil
}

func (gp *gzipWriterPool) Put(gzipWriter *GzipWriter) error {
	// Reset the writer to its initial state
	err := gzipWriter.Writer.Flush()
	if err != nil {
		return err
	}

	gzipWriter.Writer.Reset(io.Discard)

	// Put the writer back into the pool
	gp.p.Put(gzipWriter)
	return nil
}

type GzipWriter struct {
	ResWriter       http.ResponseWriter
	Writer          *gzip.Writer
	GzipWriterMutex sync.Mutex
	GzipFlag        bool
}

func (w *GzipWriter) Header() http.Header {
	return w.ResWriter.Header()
}

func (w *GzipWriter) WriteHeader(statusCode int) {
	if !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "application/json") &&
		!strings.Contains(w.ResWriter.Header().Get("Content-Type"), "text/html") {
		w.ResWriter.WriteHeader(statusCode)
		w.GzipFlag = false
		return
	}
	w.GzipFlag = true
	w.ResWriter.Header().Set("Content-Encoding", "gzip")
	w.ResWriter.WriteHeader(statusCode)
}

func (w *GzipWriter) Write(b []byte) (int, error) {
	if !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "application/json") &&
		!strings.Contains(w.ResWriter.Header().Get("Content-Type"), "text/html") {
		return w.ResWriter.Write(b)
	}
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
				// If gzip is not supported then return uncompressed page.
				next.ServeHTTP(w, r)
				return
			}

			log.Info("gzip is supported")

			gz, err := gzipWrPool.Get(w, compressorLevel)
			if err != nil {
				fmt.Println("1111111111111111111")
				return
			}
			next.ServeHTTP(gz, r)
			if gz.GzipFlag {
				err := gzipWrPool.Put(gz)
				if err != nil {
					log.Error("failed to close gzip")
					_, err := io.WriteString(w, err.Error())
					if err != nil {
						log.Error("failed to inform user")
						return
					}
					return
				}
			}
		}
		return http.HandlerFunc(fn)
	}
}
