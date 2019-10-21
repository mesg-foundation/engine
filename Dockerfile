# base Go image version.
FROM golang:1.13.0-stretch AS build

WORKDIR /project

# install dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG version
RUN go build -mod=readonly -o ./bin/engine -ldflags="-X 'github.com/mesg-foundation/engine/version.Version=$version'" core/main.go

FROM ubuntu:18.04
RUN apt-get update && \
      apt-get install -y --no-install-recommends ca-certificates=20180409 && \
      apt-get clean && \
      rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /project/bin/engine .
CMD ["./bin/engine"]
