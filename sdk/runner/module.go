package runnersdk

import (
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ModuleName is the name of this module.
const ModuleName = "runner"

// NewModule returns the module of this sdk.
func NewModule(k *Keeper) module.AppModule {
	return cosmos.NewAppModule(cosmos.NewAppModuleBasic(ModuleName), handler(k), querier(k))
}

func handler(k *Keeper) cosmos.Handler {
	return func(request cosmostypes.Request, msg cosmostypes.Msg) (hash.Hash, error) {
		switch msg := msg.(type) {
		case msgCreateRunner:
			run, err := k.Create(request, &msg)
			if err != nil {
				return nil, cosmos.NewMesgWrapError(cosmos.CodeInternal, err)
			}
			return run.Hash, nil
		case msgDeleteRunner:
			if err := k.Delete(request, &msg); err != nil {
				return nil, cosmos.NewMesgWrapError(cosmos.CodeInternal, err)
			}
			return nil, nil
		default:
			errmsg := fmt.Sprintf("Unrecognized runner Msg type: %v", msg.Type())
			return nil, cosmostypes.ErrUnknownRequest(errmsg)
		}
	}
}

func querier(k *Keeper) cosmos.Querier {
	return func(request cosmostypes.Request, path []string, req abci.RequestQuery) (res interface{}, err error) {
		switch path[0] {
		case "get":
			hash, err := hash.Decode(path[1])
			if err != nil {
				return nil, err
			}
			return k.Get(request, hash)
		case "list":
			return k.List(request)
		default:
			return nil, fmt.Errorf("unknown runner query endpoint %s", path[0])
		}
	}
}
