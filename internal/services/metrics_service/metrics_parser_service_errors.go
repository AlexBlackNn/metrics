package metrics_service

import "errors"

var (
	ErrNotValidUrl         = errors.New("not valid URL")
	ErrNotValidMetricValue = errors.New("not valid metric value")
	ErrNotValidMetricType  = errors.New("not valid metric type")
)
