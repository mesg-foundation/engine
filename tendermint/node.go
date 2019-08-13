package tendermint

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

// NewNode retruns new tendermint node that runs the app.
func NewNode(logger log.Logger, app abci.Application, cfg *config.Config, validatorPubKey crypto.PubKey) (*node.Node, error) {
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	me := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())

	return node.NewNode(cfg,
		me,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		genesisLoader(validatorPubKey),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		logger,
	)
}

// genesisLoader returns a generator for genesis structure.
func genesisLoader(validator crypto.PubKey) func() (*types.GenesisDoc, error) {
	return func() (*types.GenesisDoc, error) {
		genesis := &types.GenesisDoc{
			GenesisTime:     time.Date(2019, 8, 8, 0, 0, 0, 0, time.UTC),
			ChainID:         "mesg-chain",
			ConsensusParams: types.DefaultConsensusParams(),
			Validators: []types.GenesisValidator{{
				Address: validator.Address(),
				PubKey:  validator,
				Power:   1,
				Name:    "validator",
			}},
			AppState: []byte("{}"),
		}
		if err := genesis.ValidateAndComplete(); err != nil {
			return nil, err
		}
		return genesis, nil
	}
}