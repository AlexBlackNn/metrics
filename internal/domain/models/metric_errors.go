package models

import "errors"

var (
	ErrNotValidMetricValue    = errors.New("invalid metric value")
	ErrNotValidMetricType     = errors.New("invalid metric type")
	ErrAddDifferentMetricType = errors.New("different metric types")
	ErrAddDifferentMetricName = errors.New("different metric names")
	ErrAddMetricValueCast     = errors.New("cannot cast metric to required type")
)
