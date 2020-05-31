package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/x/process/internal/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/process/get/{hash}",
		queryGetHandlerFn(cliCtx),
	).Methods(http.MethodGet)

	r.HandleFunc(
		"/process/list",
		queryListHandlerFn(cliCtx),
	).Methods(http.MethodGet)

	r.HandleFunc(
		"/process/hash",
		queryHashHandlerFn(cliCtx),
	).Methods(http.MethodPost)

	r.HandleFunc(
		"/process/exist/{hash}",
		queryExistHandlerFn(cliCtx),
	).Methods(http.MethodGet)
}

func queryGetHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGet, vars["hash"])

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryListHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryList)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HashRequest is the request of the hash endpoint.
type HashRequest struct {
	Name  string                  `json:"name,omitempty"`
	Nodes []*process.Process_Node `json:"nodes,omitempty"`
	Edges []*process.Process_Edge `json:"edges,omitempty"`
}

func queryHashHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var req HashRequest
		if err := cliCtx.Codec.UnmarshalJSON(data, &req); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proc, err := process.New(req.Name, req.Nodes, req.Edges, sdk.AccAddress{})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, proc.Hash.String())
	}
}

func queryExistHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryExist, vars["hash"])

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
