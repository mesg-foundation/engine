.PHONY: all check-version docker-publish docker-publish-dev docker-tools dev lint dep build test mock protobuf changelog clean genesis

MAJOR_VERSION := $(shell echo $(version) | cut -d . -f 1)	
MINOR_VERSION := $(shell echo $(version) | cut -d . -f 1-2)
PATCH_VERSION := $(version)

all: clean lint test build e2e

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

docker-dev:
	./scripts/build-engine.sh

docker-publish: docker-build
	docker push mesg/engine:$(MAJOR_VERSION)
	docker push mesg/engine:$(MINOR_VERSION)
	docker push mesg/engine:$(PATCH_VERSION)
	docker push mesg/engine:latest

docker-publish-dev: check-version
	docker build -t mesg/engine:dev --build-arg version=$(version) .
	docker push mesg/engine:dev

docker-test:
	docker build -t mesg/test:local -f Dockerfile.e2e .

docker-tools:
	docker build -t mesg/tools:local -f Dockerfile.tools .

dev: docker-dev
	- ./scripts/dev.sh

dep:
	go mod download

build: check-version dep
	go build -mod=readonly -o ./bin/engine -ldflags="-X 'github.com/mesg-foundation/engine/version.Version=$(version)'" main.go

e2e: docker-test
	docker run \
		--mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
		mesg/test:local


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
			mesg/engine:$(version) \
			mesg/engine:latest \
			mesg/engine:dev 2>/dev/null

genesis:
	go run internal/tools/gen-genesis/main.go --path $(path)  --chain-id $(chain-id) --validators $(validators)
