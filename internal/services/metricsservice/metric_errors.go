package metricsservice

import "errors"

var (
	ErrNotValidURL          = errors.New("not valid URL")
	ErrCouldNotUpdateMetric = errors.New("could not update metric")
	ErrMetricNotFound       = errors.New("metric not found")
	ErrCouldNotGetMetric    = errors.New("could not get metric")
)
