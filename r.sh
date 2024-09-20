#!/bin/bash


go build -o cmd/agent/agent cmd/agent/main.go
go build -o cmd/server/server cmd/server/main.go
./metricstest -test.v -test.run=^TestIteration7 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=.
