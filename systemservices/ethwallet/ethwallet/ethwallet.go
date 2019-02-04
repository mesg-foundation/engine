package ethwallet

import (
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/mesg-foundation/core/client/service"
)

const (
	keystoreEnv = "MESG_KEYSTORE"
)

// Ethwallet is a Ethereum Wallet
type Ethwallet struct {
	service  *service.Service
	keystore *keystore.KeyStore
}

// New creates a new instance of EthWallet
func New() (*Ethwallet, error) {
	service, err := service.New()
	if err != nil {
		return nil, err
	}

	ks := keystore.NewKeyStore(os.Getenv(keystoreEnv), keystore.StandardScryptN, keystore.StandardScryptP)

	return &Ethwallet{
		service:  service,
		keystore: ks,
	}, nil
}

// Listen listens for tasks from MESG
func (ethwallet *Ethwallet) Listen() error {
	return ethwallet.service.Listen(
		service.Task("list", ethwallet.list),
		service.Task("create", ethwallet.create),
		service.Task("delete", ethwallet.delete),
		service.Task("export", ethwallet.export),
		service.Task("import", ethwallet.importA),
		service.Task("sign", ethwallet.sign),
	)
}
