#!/bin/bash

go-bindata -o service/importer/assets/schema.go -pkg assets ./service/importer/assets/schema.json
go-bindata -o cmd/service/assets/readmeTemplate.go -pkg assets ./cmd/service/assets/readmeTemplate.md