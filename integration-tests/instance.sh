#!/bin/bash -e

# Extract service hash
serviceHash=$(jq -r .hash ./service-create-response.json)

# Create instance
# TODO: add envs
jq -n --arg serviceHash $serviceHash '{"serviceHash":$serviceHash}' > instance-get-request.json
grpcurl -plaintext -d "$(cat instance-get-request.json)" localhost:50052 api.Instance/Create > instance-create-response.json

# Extract instance hash
instanceHash=$(jq -r .hash ./instance-create-response.json)

# Get instance
jq -n --arg hash $instanceHash '{"hash":$hash}' > instance-get-request.json
grpcurl -plaintext -d "$(cat instance-get-request.json)" localhost:50052 api.Instance/Get > instance-get-response.json

# List instances
# TODO: add serviceHash filter
jq -n '{}' > instance-list-request.json
grpcurl -plaintext -d "$(cat instance-list-request.json)" localhost:50052 api.Instance/List > instance-list-response.json

# Delete instance
jq -n --arg hash $instanceHash '{"hash":$hash, "deleteData": true}' > instance-delete-request.json
grpcurl -plaintext -d "$(cat instance-delete-request.json)" localhost:50052 api.Instance/Delete > instance-delete-response.json
