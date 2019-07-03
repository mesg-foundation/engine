package ethwallet

import (
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/client"
)

const (
	keystoreEnv = "MESG_KEYSTORE"
)

// Ethwallet is a Ethereum Wallet
type Ethwallet struct {
	client   *client.Client
	keystore *keystore.KeyStore
}

// New creates a new instance of EthWallet
func New() (*Ethwallet, error) {
	client, err := client.New()
	if err != nil {
		return nil, err
	}

	keystore := keystore.NewKeyStore(os.Getenv(keystoreEnv), keystore.StandardScryptN, keystore.StandardScryptP)
	return &Ethwallet{
		client:   client,
		keystore: keystore,
	}, nil
}

// Listen listens for tasks from MESG
func (e *Ethwallet) Listen() error {
	t := e.client.TaskRunner()
	t.Add("list", e.list)
	t.Add("create", e.create)
	t.Add("delete", e.delete)
	t.Add("export", e.export)
	t.Add("import", e.importA)
	t.Add("sign", e.sign)
	t.Add("importFromPrivateKey", e.importFromPrivateKey)
	return t.Run()
}
