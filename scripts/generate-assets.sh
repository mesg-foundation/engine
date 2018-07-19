#!/bin/bash

go-bindata -o service/serialize/assets/schema.go -pkg assets ./service/serialize/assets/schema.json
go-bindata -o cmd/service/assets/readmeTemplate.go -pkg assets ./cmd/service/assets/readmeTemplate.md