FROM golang:1.11.4
WORKDIR /project
COPY . .
RUN go install -v ./...
RUN cd /go/bin
CMD ["cmd"]
