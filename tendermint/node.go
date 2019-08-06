package tendermint

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

// NewNode retruns new tendermint node that runs the app.
func NewNode(logger log.Logger, app abci.Application, root string, seeds string) (*node.Node, error) {
	cfg := config.DefaultConfig()
	cfg.P2P.Seeds = seeds
	cfg.SetRoot(root)

	var (
		validator = privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
		nodeKey   = &p2p.NodeKey{
			PrivKey: ed25519.GenPrivKey(),
		}
	)

	return node.NewNode(cfg,
		validator,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		logger,
	)
}
