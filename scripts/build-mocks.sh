#!/bin/bash -e

# generate mocks
mockery -name ExecutionSDK -dir ./orchestrator -output ./orchestrator/mocks
mockery -name ResultSDK -dir ./orchestrator -output ./orchestrator/mocks
mockery -name EventSDK -dir ./orchestrator -output ./orchestrator/mocks
mockery -name ProcessSDK -dir ./orchestrator -output ./orchestrator/mocks
mockery -name RunnerSDK -dir ./orchestrator -output ./orchestrator/mocks
