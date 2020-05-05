# base Go image version for building the binaries
FROM golang:1.13.10 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG version

RUN go build -mod=readonly -o ./bin/mesg-cli -ldflags="-s -w -X 'github.com/cosmos/cosmos-sdk/version.Name=mesg' -X 'github.com/cosmos/cosmos-sdk/version.ServerName=mesg-daemon' -X 'github.com/cosmos/cosmos-sdk/version.ClientName=mesg-cli' -X 'github.com/cosmos/cosmos-sdk/version.Version=$version'" ./cmd/mesg-cli/
RUN go build -mod=readonly -o ./bin/mesg-daemon -ldflags="-s -w -X 'github.com/cosmos/cosmos-sdk/version.Name=mesg' -X 'github.com/cosmos/cosmos-sdk/version.ServerName=mesg-daemon' -X 'github.com/cosmos/cosmos-sdk/version.ClientName=mesg-cli' -X 'github.com/cosmos/cosmos-sdk/version.Version=$version'" ./cmd/mesg-daemon/

# ubuntu image with binaries for distribution
FROM ubuntu:18.04
RUN apt-get update && \
  apt-get install -y --no-install-recommends ca-certificates=20180409 && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /app
ENV PATH="/app:${PATH}"

COPY --from=build /app/bin/mesg-cli .
COPY --from=build /app/bin/mesg-daemon .

CMD ["mesg-daemon", "start"]
