#!/bin/bash -e

# Create Service
grpcurl -plaintext -d "$(cat service-definition.json)" localhost:50052 api.ServiceX/Create > service-create-response.json

# Extract service hash
serviceHash=$(jq -r .hash ./service-create-response.json)

# Get service
jq -n --arg hash $serviceHash '{"hash":$hash}' > service-get-request.json
grpcurl -plaintext -d "$(cat service-get-request.json)" localhost:50052 api.ServiceX/Get > service-get-response.json

# List services
jq -n '{}' > service-list-request.json
grpcurl -plaintext -d "$(cat service-list-request.json)" localhost:50052 api.ServiceX/List > service-list-response.json

# Delete service
jq -n --arg hash $serviceHash '{"hash":$hash}' > service-delete-request.json
grpcurl -plaintext -d "$(cat service-delete-request.json)" localhost:50052 api.ServiceX/Delete > service-delete-response.json
