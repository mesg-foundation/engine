package cosmos

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

// RegisterSimulateRoute registers the route on the router.
func RegisterSimulateRoute(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/txs/simulate", SimulateRequestHandlerFn(cliCtx)).Methods("POST")
}

// SimulateReq defines the properties of a send request's body.
type SimulateReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Msgs    []sdk.Msg    `json:"msgs" yaml:"msgs"`
}

// SimulateRequestHandlerFn - http request handler to simulate msgs.
func SimulateRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SimulateReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// force simulate to true to use the following function in simulate only mode.
		req.BaseReq.Simulate = true
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, req.Msgs)
	}
}
