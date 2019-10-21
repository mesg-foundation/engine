.PHONY: all docker docker-dev docker-tools dev dep build test mock protobuf changelog clean 

all: clean test build

docker-dev:
	docker build -t mesg/engine:dev --build-arg version=dev .

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
	./scripts/changelog.sh $(milestone)

clean:
	- rm -rf bin/*
	- docker image rm mesg/engine:dev
