#!/bin/bash
./metricstest -test.v -test.run=^TestIteration1$ -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server
./metricstest -test.v -test.run=^TestIteration2 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server  -source-path=.
./metricstest -test.v -test.run=^TestIteration3 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=.
./metricstest -test.v -test.run=^TestIteration4 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. --server-port 8000
./metricstest -test.v -test.run=^TestIteration5 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. --server-port 8001
