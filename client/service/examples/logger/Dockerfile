FROM golang:1.11.4
WORKDIR /project
COPY . .
RUN go build -o service-logger ./cmd
CMD ["./cmd"]
