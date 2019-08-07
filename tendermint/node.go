package tendermint

import (
	"os"

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
func NewNode(logger log.Logger, app abci.Application, root string, seeds string, validator string) (*node.Node, error) {

	os.MkdirAll(root+"/config", os.FileMode(0755))
	cfg := config.DefaultConfig()
	cfg.P2P.Seeds = seeds
	cfg.SetRoot(root)

	os.MkdirAll(root+"/data", os.FileMode(0755))
	var me = privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
	nodeKey := &p2p.NodeKey{
		PrivKey: ed25519.GenPrivKey(),
	}

	validatorPubKey := me.GetPubKey()
	if validator != "" {
		// TOFIX: this is not working
		// TODO: convert string to pubkey
		var pubTmp ed25519.PubKeyEd25519
		copy(pubTmp[:], validator)
		validatorPubKey = pubTmp
		logger.Error("will use validator", validatorPubKey)
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
