ARG from=golang:1.13.10
FROM $from AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG version

RUN make build

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
