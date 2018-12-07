#!/bin/bash

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

go-bindata -o service/importer/assets/schema.go -pkg assets ./service/importer/assets/schema.json
go-bindata -o commands/provider/assets/readme_template.go -pkg assets ./commands/provider/assets/readme_template.md
