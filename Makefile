.PHONY: all docker docker-dev docker-tools dev dep build test mock protobuf changelog clean 

all: clean test build docker

MAJOR_VERSION := $(shell git describe --abbrev=0 | cut -d . -f 1)
MINOR_VERSION := $(shell git describe --abbrev=0 | cut -d . -f 1-2)
PATCH_VERSION := $(shell git describe --abbrev=0)
DEV_VERSION := $(shell git describe)

docker:
	docker build \
		--build-arg version=$(PATCH_VERSION) \
		-t mesg/engine:$(MAJOR_VERSION) \
		-t mesg/engine:$(MINOR_VERSION) \
		-t mesg/engine:$(PATCH_VERSION) \
		-t mesg/engine:latest \
		.

docker-dev:
	docker build -t mesg/engine:dev --build-arg version=$(DEV_VERSION) -f Dockerfile .

docker-tools:
	docker build -t mesg/tools:local -f Dockerfile.tools .

dev: docker-dev
	- ./scripts/dev.sh

dep:
	go mod download

build: dep
	go build -mod=readonly -o ./bin/engine -ldflags="-X 'github.com/mesg-foundation/engine/version.Version=$(PATCH_VERSION)'" core/main.go

test:
	go test ./...

mock: docker-tools
	docker run --rm -v $(PWD):/project mesg/tools:local	./scripts/build-mocks.sh

protobuf: docker-tools
	docker run --rm -v $(PWD):/project mesg/tools:local	./scripts/build-proto.sh

changelog:
	./scripts/changelog.sh

clean:
	- rm -rf bin/*
	- docker image rm \
			mesg/engine:$(MAJOR_VERSION) \
			mesg/engine:$(MINOR_VERSION) \
			mesg/engine:$(PATCH_VERSION) \
			mesg/engine:latest
