#!/bin/bash -e

# generate mocks
mockery -name Container -dir ./container -output ./container/mocks
mockery -name CommonAPIClient -dir ./internal/mocks -output ./utils/docker/mocks
mockery -name ExecutionSDK -dir ./orchestrator -output ./orchestrator/mocks
mockery -name EventSDK -dir ./orchestrator -output ./orchestrator/mocks
mockery -name ProcessSDK -dir ./orchestrator -output ./orchestrator/mocks
