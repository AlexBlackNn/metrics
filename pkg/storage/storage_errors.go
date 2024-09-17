package storage

import "errors"

var (
	ErrMetricNotFound     = errors.New("metric not found")
	ErrSQLExec            = errors.New("sql query failed")
	ErrConnectionFailed   = errors.New("connection failed")
	ErrUnexpectedBehavior = errors.New("unexpected behavior")
)
