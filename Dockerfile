# base Go image version.
FROM golang:1.14.0-alpine3.11 AS build
WORKDIR /app

RUN apk add build-base gcc abuild binutils binutils-doc gcc-doc

# install dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG version

RUN go build -mod=readonly -o ./bin/engine -ldflags="-s -w -X 'github.com/mesg-foundation/engine/version.Version=$version'" core/main.go

FROM alpine:3.11.3
WORKDIR /app

RUN apk add --no-cache ca-certificates apache2-utils

COPY --from=build /app/bin/engine .
CMD ["./engine"]
