package gzipcompressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var gzipWrPool = &gzipWriterPool{}

type GzipWriter struct {
	ResWriter http.ResponseWriter
	Writer    *gzip.Writer
	GzipFlag  bool
}

func (w *GzipWriter) Header() http.Header {
	return w.ResWriter.Header()
}

func (w *GzipWriter) WriteHeader(statusCode int) {
	if !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "application/json") &&
		!strings.Contains(w.ResWriter.Header().Get("Content-Type"), "text/html") {
		w.GzipFlag = false
		w.ResWriter.WriteHeader(statusCode)
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

func (w *GzipWriter) Flush() error {
	return w.Writer.Flush()
}

func (w *GzipWriter) Reset() {
	w.Writer.Reset(io.Discard)
}

func GzipCompressor(compressorLevel int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if !strings.Contains(strings.Join(r.Header.Values("Accept-Encoding"), " "), "gzip") {
				// If gzip is not supported then return uncompressed page.
				next.ServeHTTP(w, r)
				return
			}

			gzipWr, err := gzipWrPool.get(w, compressorLevel)
			if err != nil {
				_, err = io.WriteString(w, err.Error())
				if err != nil {
					return
				}
				return
			}
			next.ServeHTTP(gzipWr, r)
			if gzipWr.GzipFlag {
				err = gzipWrPool.put(gzipWr)
			} else {
				err = gzipWrPool.putNoFlush(gzipWr)
			}
			if err != nil {
				_, err = io.WriteString(w, err.Error())
				if err != nil {
					return
				}
				return
			}
		}
		return http.HandlerFunc(fn)
	}
}
