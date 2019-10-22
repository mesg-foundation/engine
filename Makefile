.PHONY: all check-env docker-build docker-dev docker-tools dev lint dep build test mock protobuf changelog clean 

all: clean test build docker-build

check-env:
ifndef version
	$(error version is undefined)
endif

docker-build: check-env
	docker build -t mesg/engine:$(version) --build-arg version=$(version) .

docker-dev:
	docker build -t mesg/engine:dev --build-arg version=dev .

docker-tools:
	docker build -t mesg/tools:local -f Dockerfile.tools .

dev: docker-dev
	./scripts/dev.sh

dep:
	go mod download

build: check-env dep
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
			mesg/engine:$(version) \
			mesg/engine:latest \
			mesg/engine:dev 2>/dev/null
