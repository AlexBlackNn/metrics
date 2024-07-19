#!/bin/bash
cd infra && "docker compose" up -d
sleep 15
cd ..
go run ./cmd/server/migrator  --p ./migrations --d postgres://app:app123@localhost:5432/metric_db?sslmode=disable
