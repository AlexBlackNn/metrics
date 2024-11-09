package ipchecker

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"

	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/go-chi/chi/v5/middleware"
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

func IPChecker(log *slog.Logger, cfg *configserver.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/IPChecker"),
		)
		log.Info("IPChecker middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {

			if cfg.TrustedSubnet == "*" {
				ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
				next.ServeHTTP(ww, r)
				return
			}

			ipStr := r.Header.Get("X-Real-Ip")
			if ipStr == "" {
				ipStr = r.Header.Get("X-Forwarded-For")
			}
			_, trustedNet, err := net.ParseCIDR(cfg.TrustedSubnet)
			if err != nil {
				log.Error("parse CIDR error for trusted Subnet, check cfg:", "err", err)
				return
			}
			ip := net.ParseIP(ipStr)
			if (ip == nil) || (!trustedNet.Contains(ip)) {
				log.Warn("ip validation failed")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				dataMarshal, _ := json.Marshal(Error("access forbidden"))
				_, err := w.Write(dataMarshal)
				if err != nil {
					log.Error("write error:", "err", err)
				}
				return
			}
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
