FROM golang:1.11.4
WORKDIR /project
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build ./main.go
CMD ["./main"]
