FROM golang
RUN mkdir -p src/github.com/mesg-foundation/core
ADD . src/github.com/mesg-foundation/core
WORKDIR src/github.com/mesg-foundation/core
RUN go get -v -t ./...
RUN go build -o mesg-core cli/main.go
CMD ["./mesg-core"]
