FROM golang:1.13

WORKDIR /app

COPY . .

RUN go build -o "mesg-test"

CMD [ "mesg-test" ]
