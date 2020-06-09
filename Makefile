.PHONY: build build-multiplatform build-docker publish publish-docker-prod publish-docker-unstable test e2e dev dev-mon dep lint build-tools protobuf changelog clean

version ?= local
MAJOR_VERSION := $(shell echo $(version) | cut -d . -f 1)
MINOR_VERSION := $(shell echo $(version) | cut -d . -f 1-2)

LDFLAGS = "-s -w -X 'github.com/cosmos/cosmos-sdk/version.Name=mesg' -X 'github.com/cosmos/cosmos-sdk/version.ServerName=mesg-daemon' -X 'github.com/cosmos/cosmos-sdk/version.ClientName=mesg-cli' -X 'github.com/cosmos/cosmos-sdk/version.Version=$(version)'"
OSS = linux darwin
ARCHS = amd64 386

# Build

build: dep
	CGO_ENABLED=0 go build -mod=readonly -o "./bin/mesg-daemon" -ldflags=$(LDFLAGS) ./cmd/mesg-daemon/
	CGO_ENABLED=0 go build -mod=readonly -o "./bin/mesg-cli" -ldflags=$(LDFLAGS) ./cmd/mesg-cli/

build-multiplatform: dep
	for os in $(OSS) ; do \
		for arch in $(ARCHS) ; do \
			echo "Building $$os $$arch..." ; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -mod=readonly -o "./bin/mesg-daemon-$$os-$$arch" -ldflags=$(LDFLAGS) ./cmd/mesg-daemon/ ; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -mod=readonly -o "./bin/mesg-cli-$$os-$$arch" -ldflags=$(LDFLAGS) ./cmd/mesg-cli/ ; \
		done \
	done

build-docker-cache:
	# building cache image
	docker build \
		--build-arg version=$(version) \
		--target build \
		-t mesg/engine:$(version)-build \
		.

build-docker-cache-if-needed:
	if [ -z "$(shell docker images -q mesg/engine:$(version)-build)" ]; then \
		make build-docker-cache ; \
	fi

build-docker: build-docker-cache-if-needed
	# building image
	docker build \
		--build-arg version=$(version) \
		--build-arg from=mesg/engine:$(version)-build \
		-t mesg/engine:$(version) \
		.
	# building dev image
	docker build \
		-f ./Dockerfile.dev \
		--build-arg from=mesg/engine:$(version) \
		-t mesg/engine:$(version)-dev \
		.

# Publish

publish-unstable: build-multiplatform
	go install github.com/tcnksm/ghr
	ghr -u mesg-foundation -r engine -p 1 -delete -prerelease -n "Unstable release" -b "Warning - this is an unstable release, use it only if you know what are doing." unstable ./bin

publish-prod: build-multiplatform
	go install github.com/tcnksm/ghr
	ghr -u mesg-foundation -r engine -p 1 -replace "$(version)" ./bin

publish-docker-prod: build-docker
	docker tag mesg/engine:$(version) mesg/engine:$(MINOR_VERSION)
	docker tag mesg/engine:$(version) mesg/engine:$(MAJOR_VERSION)
	docker tag mesg/engine:$(version) mesg/engine:latest

	docker push mesg/engine:$(version)
	docker push mesg/engine:$(MINOR_VERSION)
	docker push mesg/engine:$(MAJOR_VERSION)
	docker push mesg/engine:latest

	docker tag mesg/engine:$(version)-dev mesg/engine:$(MINOR_VERSION)-dev
	docker tag mesg/engine:$(version)-dev mesg/engine:$(MAJOR_VERSION)-dev
	docker tag mesg/engine:$(version)-dev mesg/engine:latest-dev

	docker push mesg/engine:$(version)-dev
	docker push mesg/engine:$(MINOR_VERSION)-dev
	docker push mesg/engine:$(MAJOR_VERSION)-dev
	docker push mesg/engine:latest-dev

publish-docker-unstable: build-docker
	docker tag mesg/engine:$(version) mesg/engine:unstable
	docker push mesg/engine:unstable

	docker tag mesg/engine:$(version)-dev mesg/engine:unstable-dev
	docker push mesg/engine:unstable-dev

# Test

test: dep
	go test -short -mod=readonly -v -coverprofile=coverage.txt ./...

e2e: build-docker
	./scripts/run-e2e.sh "$(version)"

dev: build-docker
	./scripts/dev.sh "$(version)"

dev-mon: build-docker
	./scripts/dev.sh "$(version)" "monitoring"

# MISC

dep:
	go mod download

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.21 golangci-lint run

build-tools:
	docker build -t mesg/tools:local -f Dockerfile.tools .

protobuf: build-tools
	docker run --rm -v $(PWD):/project mesg/tools:local	./scripts/build-proto.sh

changelog:
	./scripts/changelog.sh $(milestone)

clean:
	- rm -rf bin
	- docker volume rm engine
	- docker image rm $(shell docker images -q mesg/engine)
	- docker image rm $(shell docker images -q mesg/tools)
