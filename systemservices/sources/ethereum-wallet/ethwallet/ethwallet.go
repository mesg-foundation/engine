package ethwallet

import (
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/mesg-foundation/core/client/service"
)

const (
	keystoreEnv = "MESG_KEYSTORE"
)

type Ethwallet struct {
	service  *service.Service
	keystore *keystore.KeyStore
}

func New() (*Ethwallet, error) {
	// mesg client
	service, err := service.New()
	if err != nil {
		return nil, err
	}

	// keystore
	ks := keystore.NewKeyStore(os.Getenv(keystoreEnv), keystore.StandardScryptN, keystore.StandardScryptP)

	return &Ethwallet{
		service:  service,
		keystore: ks,
	}, nil
}

func (ethwallet *Ethwallet) Listen() error {
	return ethwallet.service.Listen(
		service.Task("list", ethwallet.list),
		service.Task("new", ethwallet.new),
		service.Task("delete", ethwallet.delete),
		service.Task("export", ethwallet.export),
		service.Task("import", ethwallet.importA),
		service.Task("sign", ethwallet.sign),
	)
}
