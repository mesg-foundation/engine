FROM golang:1.13.8-alpine3.11
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o "mesg-test" -ldflags="-s -w" ./main.go
CMD [ "./mesg-test" ]
