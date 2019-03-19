FROM golang:1.12.0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main

CMD ["./main"]
