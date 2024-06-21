package v2

import (
	"encoding/json"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/AlexBlackNn/metrics/internal/lib/response"
	"net/http"
)

func responseOK(w http.ResponseWriter, r *http.Request, metric models.MetricInteraction) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if metric.GetType() == "counter" {
		metricValue := int64(metric.GetValue().(uint64))
		metricMarshal, _ := json.Marshal(Metrics{
			ID:    metric.GetName(),
			MType: metric.GetType(),
			Delta: &metricValue,
		})
		w.Write(metricMarshal)
		return
	}
	metricValue := metric.GetValue().(float64)
	metricMarshal, _ := json.Marshal(Metrics{
		ID:    metric.GetName(),
		MType: metric.GetType(),
		Value: &metricValue,
	})
	w.Write(metricMarshal)

}

func responseError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errMarshal, _ := json.Marshal(response.Error(message))
	w.Write(errMarshal)
}
