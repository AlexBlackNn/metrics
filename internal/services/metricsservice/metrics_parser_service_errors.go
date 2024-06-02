package metricsservice

import "errors"

var (
	ErrNotValidURL          = errors.New("not valid URL")
	ErrNotValidMetricValue  = errors.New("not valid metric value")
	ErrNotValidMetricType   = errors.New("not valid metric type")
	ErrCouldNotUpdateMetric = errors.New("could not update metric")
)
