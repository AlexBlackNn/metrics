package decryptor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

func Decryptor(log *slog.Logger, cfg *configserver.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/Decryptor"),
		)

		log.Info("Decryptor middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			// Проверяем, есть ли заголовок, указывающий на шифрование
			if r.Header.Get("X-Encrypted") == "true" {
				fmt.Println("11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
				// Читаем тело запроса
				encryptedBody, err := io.ReadAll(r.Body)
				if err != nil {
					log.Error("error reading request body", "err", err)
					http.Error(w, "error reading request body", http.StatusInternalServerError)
					return
				}
				decryptor, err := NewDecryptor(cfg.CryptoKeyPath)
				if err != nil {
					log.Error("error creating decryptor", "err", err)
					return
				}
				// Дешифруем сообщение
				decryptedBody, err := decryptor.DecryptMessage(string(encryptedBody))
				if err != nil {
					log.Error("error decrypting message", "err", err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					dataMarshal, _ := json.Marshal(Error("decryption failed"))
					_, err := w.Write(dataMarshal)
					if err != nil {
						log.Error("write error:", "err", err)
					}
					return
				}

				// Устанавливаем расшифрованное тело обратно в запрос
				r.Body = io.NopCloser(bytes.NewBuffer([]byte(decryptedBody)))
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			entry.Info("request completed",
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
			)
		}

		return http.HandlerFunc(fn)
	}
}
