FROM golang:1.10.4
WORKDIR /go/src/github.com/mesg-foundation/core/systemservices/-sources/resolver
RUN go get -v github.com/mesg-foundation/go-service && \
    go get -v github.com/mesg-foundation/go-service/mesgtest && \
    go get -v github.com/mesg-foundation/core/x/xstrings && \
    go get -v github.com/mesg-foundation/core/interface/grpc/core && \
    go get -v google.golang.org/grpc && \
    go get -v github.com/stvp/assert
COPY . .
RUN go get -v -d ./...
RUN go install ./...
CMD ["service-system-resolver"]