package v3

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/go-playground/validator/v10"
	"github.com/mailru/easyjson"
	"net/http"
	"strings"
)

type Metrics struct {
	ID    string   `json:"id"`                                  // metrics name
	MType string   `json:"type" validate:"oneof=gauge counter"` // mType = counter || gauge
	Delta *int64   `json:"delta,omitempty"`                     // exists if mType = counter
	Value *float64 `json:"value,omitempty"`                     // exists if mType = gauge
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const StatusError = "Error"

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func HealthOk(msg string) Response {
	return Response{
		Status: msg,
	}
}

func ValidationError(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return strings.Join(errMsgs, ", ")
}

func responseOK(w http.ResponseWriter, r *http.Request, metric models.MetricGetter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if metric.GetType() == configserver.MetricTypeCounter {
		metricValue := int64(metric.GetValue().(uint64))
		metricMarshal, _ := easyjson.Marshal(Metrics{
			ID:    metric.GetName(),
			MType: metric.GetType(),
			Delta: &metricValue,
		})
		w.Write(metricMarshal)
		return
	}
	metricValue := metric.GetValue().(float64)
	metricMarshal, _ := easyjson.Marshal(Metrics{
		ID:    metric.GetName(),
		MType: metric.GetType(),
		Value: &metricValue,
	})
	w.Write(metricMarshal)
}

func responseError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	dataMarshal, _ := easyjson.Marshal(Error(message))
	w.Write(dataMarshal)
}
