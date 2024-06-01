package update

import (
	"log/slog"
	"net/http"
)

// @Summary J,y
// @Description Создает новое выражение на сервере
// @Tags Calculations
// @Accept json
// @Produce json
// @Param body body Request true "Запрос на создание выражения"
// @Success 201 {object} Response
// @Router /expression [post]
// @Security BearerAuth
func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("update", slog.String("update", "22"))
		w.Write([]byte("Привет!"))
	}
}
