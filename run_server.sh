#!/bin/bash
cd infra && "docker compose" up -d
sleep 15
cd ..
go run ./cmd/server/migrator  --migrations-path=./migrations
