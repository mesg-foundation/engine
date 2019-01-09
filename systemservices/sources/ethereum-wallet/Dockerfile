FROM golang:1.11.4
WORKDIR /go/src/github.com/mesg-foundation/core/systemservices/sources/ethereum-wallet
COPY go.mod go.sum ./
RUN export GO111MODULE=on; go mod download
COPY . .
RUN export GO111MODULE=on; go build ./main.go
CMD ["./main"]
