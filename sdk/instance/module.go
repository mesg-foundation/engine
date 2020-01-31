package instancesdk

import (
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ModuleName is the name of this module.
const ModuleName = "instance"

// NewModule returns the module of this sdk.
func NewModule(k *Keeper) module.AppModule {
	return cosmos.NewAppModule(cosmos.NewAppModuleBasic(ModuleName), handler(k), querier(k))
}

func handler(k *Keeper) cosmos.Handler {
	return func(request cosmostypes.Request, msg cosmostypes.Msg) (hash.Hash, error) {
		errmsg := fmt.Sprintf("Unrecognized instance Msg type: %v", msg.Type())
		return nil, cosmostypes.ErrUnknownRequest(errmsg)
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
			var f api.ListInstanceRequest_Filter
			if err := codec.UnmarshalBinaryBare(req.Data, &f); err != nil {
				return nil, err
			}
			return k.List(request, &f)
		case "exists":
			hash, err := hash.Decode(path[1])
			if err != nil {
				return nil, err
			}
			return k.Exists(request, hash)
		default:
			return nil, fmt.Errorf("unknown instance query endpoint %s", path[0])
		}
	}
}
