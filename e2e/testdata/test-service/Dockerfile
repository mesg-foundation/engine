FROM golang:1.14.0-alpine3.11
WORKDIR /app

RUN apk add build-base gcc abuild binutils binutils-doc gcc-doc

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -mod=readonly -o "mesg-test" -ldflags="-s -w" ./main.go
CMD [ "./mesg-test" ]
