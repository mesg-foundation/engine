package tendermint

import (
	"encoding/json"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

// NewNode retruns new tendermint node that runs the app.
func NewNode(logger log.Logger, app abci.Application, root string, seeds string, appState json.RawMessage) (*node.Node, error) {

	os.MkdirAll(root+"/config", os.FileMode(0755))
	cfg := config.DefaultConfig()
	cfg.P2P.Seeds = seeds
	cfg.SetRoot(root)

	os.MkdirAll(root+"/data", os.FileMode(0755))
	var validator = privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	return node.NewNode(cfg,
		validator,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		// genesisLoader(appState, validator),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		logger,
	)
}

func genesisLoader(appState json.RawMessage, validator types.PrivValidator) func() (*types.GenesisDoc, error) {
	return func() (*types.GenesisDoc, error) {
		genesis := &types.GenesisDoc{
			ChainID:         "xxx",
			ConsensusParams: types.DefaultConsensusParams(),
			Validators:      []types.GenesisValidator{
				// types.GenesisValidator{
				// 	Address: validator.GetPubKey().Address(),
				// 	PubKey:  validator.GetPubKey(),
				// 	Power:   1,
				// 	Name:    "validator",
				// },
			},
			AppState: appState,
		}
		if err := genesis.ValidateAndComplete(); err != nil {
			panic(err)
		}
		return genesis, nil
	}
}
