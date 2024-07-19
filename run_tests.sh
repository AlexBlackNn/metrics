#!/bin/bash

# Define the log file
LOG_FILE="test_log_file.txt"

# Run your test commands, capturing output to the log file
./metricstest -test.v -test.run=^TestIteration1$ -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server > "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration2 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration3 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration4 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. -server-port 8000 >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration5 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. -server-port 8001 >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration6 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration7 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration8 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. -server-port=8080 >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration9 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. -server-port=8080 -file-storage-path ./test.json >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration10 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. -server-port=8080 -file-storage-path ./test.json -database-dsn postgresql://app:app123@127.0.0.1:5432/metric_db?sslmode=disable>> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration11 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. -server-port=8080 -file-storage-path ./test.json -database-dsn postgresql://app:app123@127.0.0.1:5432/metric_db?sslmode=disable >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration12 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. -server-port=8080 -file-storage-path ./test.json -database-dsn postgresql://app:app123@127.0.0.1:5432/metric_db?sslmode=disable >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration13 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. -server-port=8080 -file-storage-path ./test.json -database-dsn postgresql://app:app123@127.0.0.1:5432/metric_db?sslmode=disable >> "$LOG_FILE" 2>&1
./metricstest -test.v -test.run=^TestIteration14 -agent-binary-path=cmd/agent/agent -binary-path=./cmd/server/server -source-path=. -server-port=8080 -file-storage-path ./test.json -database-dsn postgresql://app:app123@127.0.0.1:5432/metric_db?sslmode=disable -key 1 >> "$LOG_FILE" 2>&1

# Process the log file to find failed tests
while read line; do
  # Check if the line indicates a failed test
  if [[ "$line" =~ "FAIL" ]]; then
    echo "$line"
  fi
done < "$LOG_FILE"