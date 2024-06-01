package storage

import (
	"context"
	"orchestrator/internal/domain/models"
	"time"
)

type OperationStorageInterface interface {
	SaveOperation(
		ctx context.Context,
		operationModel models.Operation,
		appUser models.User,
		value any,
	) error
	GetOperation(
		ctx context.Context,
		operation string,
	) (models.Operation, error)
	GetOperationById(
		ctx context.Context,
		id string,
	) (models.Operation, error)
	UpdateOperation(
		ctx context.Context,
		operation models.Operation,
	) error
	GetOperationsFastPagination(
		ctx context.Context,
		pageSize int,
		cursor string,
	) ([]models.Operation, error)
	GetUserOperationsFastPagination(
		ctx context.Context,
		pageSize int,
		cursor string,
		appUser models.User,
	) ([]models.Operation, error)
}

type OperationType int

const (
	PlusOperation OperationType = iota
	MinusOperation
	MultiplicationOperation
	DivisionOperation
)

type SettingsStorageInterface interface {
	UpdateSettingsExecutionTime(
		ctx context.Context,
		opType OperationType,
		executionTime int,
	) error
	GetOperationExecutionTime(
		ctx context.Context,
	) (models.Settings, error)
}

type OperationCacheInterface interface {
	SaveOperation(ctx context.Context, operation string, result float64, ttl time.Duration) error
	GetOperation(ctx context.Context, operation string) (float64, error)
}
