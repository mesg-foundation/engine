.PHONY: all e2e check-version docker-publish docker-publish-dev docker-tools dev dev-stop dev-start lint dep build test mock protobuf changelog clean genesis clean-build clean-docker build-cmd-cosmos publish-cmd-cosmos-prod publish-cmd-cosmos-dev get-ghr

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

dev-start: docker-dev
	./scripts/dev.sh -q

dev-stop: docker-dev
	./scripts/dev.sh stop

dep:
	go mod download

build: check-version dep
	go build -mod=readonly -o ./bin/engine -ldflags="-X 'github.com/mesg-foundation/engine/version.Version=$(version)'" core/main.go

get-ghr:
	go get -u github.com/tcnksm/ghr

build-cmd-cosmos: dep
	go build -mod=readonly -o ./bin/mesg-cosmos ./cmd/mesg-cosmos/main.go

publish-cmd-cosmos-dev: get-ghr
	ghr -u mesg-foundation -r engine -delete -prerelease -n "Developer Release" -b "Warning - this is a developer release, use it only if you know what are doing. Make sure to pull the latest \`mesg/engine:dev\` image. \`\`\`docker pull mesg/engine:dev\`\`\`" release-dev ./bin/mesg-cosmos

publish-cmd-cosmos-prod: get-ghr
ifndef tag
	$(error tag is undefined)
endif
	ghr -u mesg-foundation -r engine -delete $(tag) ./bin/mesg-cosmos

e2e: docker-dev
	./scripts/run-e2e.sh

test: dep
	go test -short -mod=readonly -v -coverprofile=coverage.txt ./...

lint:
	golangci-lint run

mock: docker-tools
	docker run --rm -v $(PWD):/project mesg/tools:local	./scripts/build-mocks.sh

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
			mesg/engine:dev 2>/dev/null

clean: clean-build clean-docker

genesis:
	go run internal/tools/gen-genesis/main.go --path $(path) --chain-id $(chain-id) --validators $(validators)
