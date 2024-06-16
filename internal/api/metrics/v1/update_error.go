package v1

import "errors"

var (
	ErrNotValidURL          = errors.New("invalid URL")
	ErrNotValidMetricValue  = errors.New("invalid metric value")
	ErrNotValidMetricType   = errors.New("invalid metric type")
	ErrCouldNotUpdateMetric = errors.New("could not update metric")
)
