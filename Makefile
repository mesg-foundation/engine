.PHONY: all build build-cmd-cosmos changelog check-version clean clean-build clean-docker dep dev dev-mon dev-start dev-stop docker-build docker-dev docker-publish docker-publish-dev docker-tools genesis lint protobuf test

MAJOR_VERSION := $(shell echo $(version) | cut -d . -f 1)	
MINOR_VERSION := $(shell echo $(version) | cut -d . -f 1-2)
PATCH_VERSION := $(version)

all: clean lint build test e2e

check-version:
ifndef version
	$(error version is undefined)
endif

docker-build: check-version
	docker build \
		--build-arg version=$(PATCH_VERSION) \
		-t mesg/engine:$(MAJOR_VERSION) \
		-t mesg/engine:$(MINOR_VERSION) \
		-t mesg/engine:$(PATCH_VERSION) \
		-t mesg/engine:latest \
		.

docker-dev: dep
	./scripts/build-engine.sh

docker-publish: docker-build
	docker push mesg/engine:$(MAJOR_VERSION)
	docker push mesg/engine:$(MINOR_VERSION)
	docker push mesg/engine:$(PATCH_VERSION)
	docker push mesg/engine:latest

docker-publish-dev: check-version
	docker build -t mesg/engine:dev --build-arg version=$(version) .
	docker push mesg/engine:dev

docker-tools:
	docker build -t mesg/tools:local -f Dockerfile.tools .

dev: docker-dev
	- ./scripts/dev.sh

dev-mon: docker-dev
	- ./scripts/dev.sh -m

dev-start: docker-dev
	./scripts/dev.sh -q

dev-stop: docker-dev
	./scripts/dev.sh -m stop

dep:
	go mod download

build: check-version dep
	go build -mod=readonly -o ./bin/engine -ldflags="-s -w -X 'github.com/mesg-foundation/engine/version.Version=$(version)'" core/main.go

build-cmd: dep
	go build -mod=readonly -o ./bin/mesg-cli ./cmd/mesg-cli/
	go build -mod=readonly -o ./bin/mesg-daemon ./cmd/mesg-daemon/

e2e: docker-dev
	./scripts/run-e2e.sh

test: dep
	go test -short -mod=readonly -v -coverprofile=coverage.txt ./...

lint:
	golangci-lint run

protobuf: docker-tools
	docker run --rm -v $(PWD):/project mesg/tools:local	./scripts/build-proto.sh

changelog:
	./scripts/changelog.sh $(milestone)

clean-build:
	- rm -rf bin

clean-docker:
	- docker image rm \
			mesg/engine:$(version) \
			mesg/engine:latest \
			mesg/engine:local \
			mesg/engine:dev 2>/dev/null

clean: clean-build clean-docker

genesis:
	go run internal/tools/gen-genesis/main.go --path $(path) --chain-id $(chain-id) --validators $(validators)
