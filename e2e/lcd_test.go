package main

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	engineAccount *auth.BaseAccount
)

const (
	lcdEndpoint        = "http://127.0.0.1:1317/"
	lcdPostContentType = "application/json"
	pollingInterval    = 500 * time.Millisecond
	pollingTimeout     = 30 * time.Second
)

func lcdGet(t *testing.T, path string, ptr interface{}) {
	resp, err := http.Get(lcdEndpoint + path)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300, "request status code is not 2XX")
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	cosResp := rest.ResponseWithHeight{}
	require.NoError(t, cdc.UnmarshalJSON(body, &cosResp))
	if len(cosResp.Result) > 0 {
		require.NoError(t, cdc.UnmarshalJSON(cosResp.Result, ptr))
	}
}

func lcdPost(t *testing.T, path string, req interface{}, ptr interface{}) {
	cosResp := rest.ResponseWithHeight{}
	lcdPostBare(t, path, req, &cosResp)
	if len(cosResp.Result) > 0 {
		require.NoError(t, cdc.UnmarshalJSON(cosResp.Result, ptr))
	}
}

func lcdPostBare(t *testing.T, path string, req interface{}, ptr interface{}) {
	reqBody, err := cdc.MarshalJSON(req)
	require.NoError(t, err)
	resp, err := http.Post(lcdEndpoint+path, lcdPostContentType, bytes.NewReader(reqBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300, "request status code is not 2XX")
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, cdc.UnmarshalJSON(body, ptr))
}

func pollExecution(t *testing.T, executionHash hash.Hash, status execution.Status) *execution.Execution {
	timeout := time.After(pollingTimeout)
	for {
		var exec *execution.Execution
		lcdGet(t, "execution/get/"+executionHash.String(), &exec)
		if exec.Status == status {
			return exec
		}
		select {
		case <-time.After(pollingInterval):
			continue
		case <-timeout:
			t.Errorf("pollExecution timeout with execution hash %q", executionHash)
		}
	}
}

// TODO: add blacklist of execution hash or add blockHeight
func pollExecutionOfProcess(t *testing.T, processHash hash.Hash, status execution.Status, nodeKey string) *execution.Execution {
	timeout := time.After(pollingTimeout)
	for {
		var execs []*execution.Execution
		lcdGet(t, "execution/list", &execs)
		for _, exec := range execs {
			if exec.ProcessHash.Equal(processHash) && exec.Status == status && exec.NodeKey == nodeKey {
				return exec
			}
		}
		select {
		case <-time.After(pollingInterval):
			continue
		case <-timeout:
			t.Errorf("pollExecutionOfProcess timeout with process hash %q and status %q and nodeKey %q", processHash, status, nodeKey)
		}
	}
}

func lcdBroadcastMsg(t *testing.T, msg sdk.Msg) []byte {
	return lcdBroadcastMsgs(t, []sdk.Msg{msg})
}
func lcdBroadcastMsgs(t *testing.T, msgs []sdk.Msg) []byte {
	tx := createAndSignTx(t, msgs)
	req := authrest.BroadcastReq{
		Tx:   tx,
		Mode: "block", // TODO: should be changed to "sync" and wait for the tx event
	}
	var res sdk.TxResponse
	lcdPostBare(t, "txs", req, &res)
	require.Equal(t, abci.CodeTypeOK, res.Code, "transaction returned with invalid code %d: %s", res.Code, res)
	result, err := hex.DecodeString(res.Data)
	require.NoError(t, err)
	return result
}

var getAccountMutex sync.Mutex

func getAccount(t *testing.T) *auth.BaseAccount {
	getAccountMutex.Lock()
	defer getAccountMutex.Unlock()
	if engineAccount == nil {
		accKb, err := kb.GetByAddress(engineAddress)
		require.NoError(t, err)
		engineAccount = auth.NewBaseAccount(accKb.GetAddress(), nil, accKb.GetPubKey(), 0, 0)
	}
	localSeq := engineAccount.GetSequence()
	lcdGet(t, "auth/accounts/"+engineAddress.String(), &engineAccount)
	// replace seq if sup
	if localSeq > engineAccount.GetSequence() {
		engineAccount.SetSequence(localSeq)
	}
	return engineAccount
}

func createAndSignTx(t *testing.T, msgs []sdk.Msg) authtypes.StdTx {
	// retrieve account
	accR := getAccount(t)
	sequence := accR.GetSequence()
	accR.SetSequence(accR.GetSequence() + 1)

	minGasPrices, err := sdk.ParseDecCoins(cfg.Cosmos.MinGasPrices)
	require.NoError(t, err)

	// Create TxBuilder
	txBuilder := authtypes.NewTxBuilder(
		authutils.GetTxEncoder(cdc),
		accR.GetAccountNumber(),
		sequence,
		flags.DefaultGasLimit*10,
		flags.DefaultGasAdjustment,
		true,
		cfg.DevGenesis.ChainID,
		"",
		nil,
		minGasPrices,
	).WithKeybase(kb)

	// TODO: to put back
	// calculate gas
	// if txBuilder.SimulateAndExecute() {
	// 	txBytes, err := txBuilder.BuildTxForSim(msgs)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	_, adjusted, err := authutils.CalculateGas(c.QueryWithData, c.cdc, txBytes, txBuilder.GasAdjustment())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	txBuilder = txBuilder.WithGas(adjusted)
	// }

	// create StdSignMsg
	stdSignMsg, err := txBuilder.BuildSignMsg(msgs)
	require.NoError(t, err)

	// create StdTx
	stdTx := authtypes.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo)

	// sign StdTx
	signedTx, err := txBuilder.SignStdTx(cfg.Account.Name, cfg.Account.Password, stdTx, false)
	require.NoError(t, err)

	return signedTx
}
