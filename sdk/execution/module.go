package executionsdk

import (
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ModuleName is the name of this module.
const ModuleName = "execution"

// NewModule returns the module of this sdk.
func NewModule(k *Keeper) module.AppModule {
	return cosmos.NewAppModule(cosmos.NewAppModuleBasic(ModuleName), handler(k), querier(k))
}

func handler(k *Keeper) cosmos.Handler {
	return func(request cosmostypes.Request, msg cosmostypes.Msg) (hash.Hash, error) {
		switch msg := msg.(type) {
		case msgCreateExecution:
			exec, err := k.Create(request, msg)
			if err != nil {
				return nil, err
			}
			return exec.Hash, nil
		case msgUpdateExecution:
			exec, err := k.Update(request, msg)
			if err != nil {
				return nil, err
			}
			return exec.Hash, nil
		default:
			errmsg := fmt.Sprintf("Unrecognized execution Msg type: %v", msg.Type())
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
			return nil, fmt.Errorf("unknown execution query endpoint %s", path[0])
		}
	}
}
