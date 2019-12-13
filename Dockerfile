# base Go image version.
FROM golang:1.13.5-alpine AS build

WORKDIR /project

# install dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG version
RUN go build -mod=readonly -o ./bin/engine -ldflags="-s -w -X 'github.com/mesg-foundation/engine/version.Version=$version'" core/main.go

FROM alpine:3.10.3
RUN apk add --no-cache ca-certificates apache2-utils
WORKDIR /app
COPY --from=build /project/bin/engine .
CMD ["./engine"]
