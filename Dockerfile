FROM golang
RUN mkdir -p src/github.com/mesg-foundation/core
RUN go get github.com/xeipuuv/gojsonschema && \
    go get gopkg.in/yaml.v2 && \
    go get github.com/stvp/assert && \
    go get github.com/logrusorgru/aurora && \
    go get github.com/fsouza/go-dockerclient && \
    go get github.com/docker/docker/api/types/swarm && \
    go get github.com/docker/docker/api/types/mount && \
    go get github.com/spf13/viper && \
    go get github.com/spf13/cobra && \
    go get github.com/mitchellh/go-homedir && \
    go get gopkg.in/AlecAivazis/survey.v1 && \
    go get github.com/kyokomi/emoji && \
    go get github.com/briandowns/spinner && \
    go get github.com/ethereum/go-ethereum/accounts && \
    go get github.com/ethereum/go-ethereum/core/types && \
    go get github.com/golang/protobuf/proto && \
    go get golang.org/x/net/context && \
    go get google.golang.org/grpc && \
    go get github.com/cpuguy83/go-md2man
WORKDIR src/github.com/mesg-foundation/core
ADD . .
RUN go get ./...
RUN go build -o mesg-daemon daemon/start/main.go
CMD ["./mesg-daemon"]
