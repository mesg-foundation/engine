package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/mesg-foundation/engine/x/instance/internal/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/ownership/get",
		queryHandlerFn(cliCtx, types.QueryGetInstance),
	).Methods(http.MethodGet)
	r.HandleFunc(
		"/ownership/list",
		queryHandlerFn(cliCtx, types.QueryListInstances),
	).Methods(http.MethodGet)
	r.HandleFunc(
		"/instance/parameters",
		queryHandlerFn(cliCtx, types.QuerierRoute),
	).Methods("GET")
}

func queryHandlerFn(cliCtx context.CLIContext, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, path)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
