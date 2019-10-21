.PHONY: all check-env docker docker-dev docker-tools dev lint dep build test mock protobuf changelog clean 

all: clean test build docker

MAJOR_VERSION := $(shell echo $(version) | cut -d . -f 1)
MINOR_VERSION := $(shell echo $(version) | cut -d . -f 1-2)
PATCH_VERSION := $(version)

check-env:
ifndef version
	$(error version is undefined)
endif

docker: check-env
	docker build \
		--build-arg version=$(PATCH_VERSION) \
		-t mesg/engine:$(MAJOR_VERSION) \
		-t mesg/engine:$(MINOR_VERSION) \
		-t mesg/engine:$(PATCH_VERSION) \
		-t mesg/engine:latest \
		.

docker-dev:
	docker build -t mesg/engine:dev --build-arg version=$(version) .

docker-tools:
	docker build -t mesg/tools:local -f Dockerfile.tools .

dev: docker-dev
	- ./scripts/dev.sh

dep:
	go mod download

build: dep
	go build -mod=readonly -o ./bin/engine -ldflags="-X 'github.com/mesg-foundation/engine/version.Version=$(version)'" core/main.go

test: dep
	go test -mod=readonly -v -coverprofile=coverage.txt ./...

lint: dep
	golangci-lint run

mock: docker-tools
	docker run --rm -v $(PWD):/project mesg/tools:local	./scripts/build-mocks.sh

protobuf: docker-tools
	docker run --rm -v $(PWD):/project mesg/tools:local	./scripts/build-proto.sh

changelog:
	./scripts/changelog.sh $(milestone)

clean:
	- rm -rf bin/*
	- docker image rm \
			mesg/engine:$(MAJOR_VERSION) \
			mesg/engine:$(MINOR_VERSION) \
			mesg/engine:$(PATCH_VERSION) \
			mesg/engine:$(version) \
			mesg/engine:latest \
			mesg/engine:dev 2>/dev/null
