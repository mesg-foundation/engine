module github.com/mesg-foundation/engine

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/sprig v2.20.0+incompatible // indirect
	github.com/Microsoft/go-winio v0.4.12 // indirect
	github.com/Microsoft/hcsshim v0.8.6 // indirect
	github.com/btcsuite/btcd v0.0.0-20190807005414-4063feeff79a // indirect
	github.com/containerd/containerd v1.3.0 // indirect
	github.com/containerd/continuity v0.0.0-20190426062206-aaeac12a7ffc // indirect
	github.com/cosmos/cosmos-sdk v0.38.3
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/cskr/pubsub v1.0.2
	github.com/docker/cli v0.0.0-20191011045415-5d85cdacd257
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.4.2-0.20191006173954-0abbb9e4ebf1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/go-bindata/go-bindata v0.0.0-20181025070752-41975c0ccc30
	github.com/go-kit/kit v0.10.0
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.2.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.3.1 // indirect
	github.com/golang/protobuf v1.3.4
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mr-tron/base58 v1.1.3
	github.com/mwitkow/go-proto-validators v0.0.0-20190212092829-1f388280e944 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/prometheus/client_golang v1.5.1
	github.com/pseudomuto/protoc-gen-doc v1.3.1
	github.com/rcrowley/go-metrics v0.0.0-20190706150252-9beb055b7962 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	github.com/tcnksm/ghr v0.13.0
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.3
	github.com/tendermint/tm-db v0.5.0
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 // indirect
	google.golang.org/grpc v1.28.0
	gopkg.in/yaml.v2 v2.2.8
	gotest.tools v2.2.0+incompatible // indirect
)

replace github.com/tendermint/go-amino => github.com/mesg-foundation/go-amino v0.15.1-0.20200216090907-4e9dce4edc74

go 1.13
