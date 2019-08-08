package tendermint

import (
	"os"
	"path/filepath"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

// NewNode retruns new tendermint node that runs the app.
func NewNode(logger log.Logger, app abci.Application, root, seeds string, validatorPubKey crypto.PubKey) (*node.Node, error) {
	if err := os.MkdirAll(filepath.Join(root, "config"), 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(root, "data"), 0755); err != nil {
		return nil, err
	}

	cfg := config.DefaultConfig()
	cfg.P2P.Seeds = seeds
	cfg.SetRoot(root)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	me := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())

	if validatorPubKey.Equals(ed25519.PubKeyEd25519{}) {
		validatorPubKey = me.GetPubKey()
		logger.Error("Me will be validator", validatorPubKey)
	}

	return node.NewNode(cfg,
		me,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		// node.DefaultGenesisDocProviderFunc(cfg),
		genesisLoader(validatorPubKey),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		logger,
	)
}

func genesisLoader(validator crypto.PubKey) func() (*types.GenesisDoc, error) {
	return func() (*types.GenesisDoc, error) {
		genesis := &types.GenesisDoc{
			GenesisTime:     time.Date(2019, 8, 8, 0, 0, 0, 0, time.UTC),
			ChainID:         "xxx",
			ConsensusParams: types.DefaultConsensusParams(),
			Validators: []types.GenesisValidator{
				types.GenesisValidator{
					Address: validator.Address(),
					PubKey:  validator,
					Power:   1,
					Name:    "validator",
				},
			},
			AppState: []byte("{}"),
		}
		if err := genesis.ValidateAndComplete(); err != nil {
			panic(err)
		}
		return genesis, nil
	}
}
