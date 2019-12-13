# base Go image version.
FROM golang:1.13.0-stretch AS build

WORKDIR /project

# install dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG version
RUN go build -mod=readonly -o ./bin/engine -ldflags="-X 'github.com/mesg-foundation/engine/version.Version=$version'" core/main.go

FROM alpine:3.10.3
RUN apk update && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build /project/bin/engine .
CMD ["./engine"]
