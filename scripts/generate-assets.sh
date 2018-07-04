#!/bin/bash

go-bindata -o service/assets/schema.go -pkg assets ./service/assets/schema.json
go-bindata -o cmd/service/assets/readmeTemplate.go -pkg assets ./cmd/service/assets/readmeTemplate.md