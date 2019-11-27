FROM golang:1.13
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o "mesg-test"
CMD [ "./mesg-test" ]
