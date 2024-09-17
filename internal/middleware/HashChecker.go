package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const StatusError = "Error"

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func HashChecker(log *slog.Logger, cfg *configserver.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/HashChecker"),
		)

		log.Info("HashChecker middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			hashBase64 := r.Header.Get("HashSHA256")
			if hashBase64 != "" {
				dst := make([]byte, base64.StdEncoding.DecodedLen(len(hashBase64)))
				n, err := base64.StdEncoding.Decode(dst, []byte(hashBase64))
				if err != nil {
					log.Error("base64 decode:", "err", err)
					return
				}
				dst = dst[:n]

				byteData, err := io.ReadAll(r.Body)

				hashCalculator := hmac.New(sha256.New, []byte(cfg.HashKey))
				hashCalculator.Write(byteData)
				hashResult := hashCalculator.Sum(nil)

				if !bytes.Equal(dst, hashResult) {
					log.Warn("hash validation failed", "err", err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnprocessableEntity)
					dataMarshal, _ := json.Marshal(Error("hash calculation failed"))
					_, err := w.Write(dataMarshal)
					if err != nil {
						log.Error("write error:", "err", err)
					}
					return
				}

				r.Body = io.NopCloser(bytes.NewBuffer(byteData))
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
