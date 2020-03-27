package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
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

func lcdGet(path string, ptr interface{}) {
	resp, err := http.Get(lcdEndpoint + path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		panic(fmt.Errorf("request status code is not 2XX, got %d", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	cosResp := rest.ResponseWithHeight{}
	if err := cdc.UnmarshalJSON(body, &cosResp); err != nil {
		panic(err)
	}
	if len(cosResp.Result) > 0 {
		if err := cdc.UnmarshalJSON(cosResp.Result, ptr); err != nil {
			panic(err)
		}
	}
}

func lcdPost(path string, req interface{}, ptr interface{}) {
	cosResp := rest.ResponseWithHeight{}
	lcdPostBare(path, req, &cosResp)
	if len(cosResp.Result) > 0 {
		if err := cdc.UnmarshalJSON(cosResp.Result, ptr); err != nil {
			panic(err)
		}
	}
}

func lcdPostBare(path string, req interface{}, ptr interface{}) {
	reqBody, err := cdc.MarshalJSON(req)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(lcdEndpoint+path, lcdPostContentType, bytes.NewReader(reqBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		panic(fmt.Errorf("request status code is not 2XX, got %d", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := cdc.UnmarshalJSON(body, ptr); err != nil {
		panic(err)
	}
}

func pollExecution(executionHash hash.Hash, status execution.Status) *execution.Execution {
	timeout := time.After(pollingTimeout)
	for {
		var exec *execution.Execution
		lcdGet("execution/get/"+executionHash.String(), &exec)
		if exec.Status == status {
			return exec
		}
		select {
		case <-time.After(pollingInterval):
			continue
		case <-timeout:
			panic(fmt.Errorf("pollExecution timeout with execution hash %q", executionHash))
		}
	}
}

// TODO: add blacklist of execution hash or add blockHeight
func pollExecutionOfProcess(processHash hash.Hash, status execution.Status, nodeKey string) *execution.Execution {
	timeout := time.After(pollingTimeout)
	for {
		var execs []*execution.Execution
		lcdGet("execution/list", &execs)
		for _, exec := range execs {
			if exec.ProcessHash.Equal(processHash) && exec.Status == status && exec.NodeKey == nodeKey {
				return exec
			}
		}
		select {
		case <-time.After(pollingInterval):
			continue
		case <-timeout:
			panic(fmt.Errorf("pollExecutionOfProcess timeout with process hash %q and status %q and nodeKey %q", processHash, status, nodeKey))
		}
	}
}

func lcdBroadcastMsg(msg sdk.Msg) []byte {
	return lcdBroadcastMsgs([]sdk.Msg{msg})
}
func lcdBroadcastMsgs(msgs []sdk.Msg) []byte {
	tx := createAndSignTx(msgs)
	req := authrest.BroadcastReq{
		Tx:   tx,
		Mode: "block", // TODO: should be changed to "sync" and wait for the tx event
	}
	var res sdk.TxResponse
	lcdPostBare("txs", req, &res)
	if abci.CodeTypeOK != res.Code {
		panic(fmt.Errorf("transaction returned with invalid code %d: %s", res.Code, res))
	}
	result, err := hex.DecodeString(res.Data)
	if err != nil {
		panic(err)
	}
	return result
}

var getAccountMutex sync.Mutex

func getAccount() *auth.BaseAccount {
	getAccountMutex.Lock()
	defer getAccountMutex.Unlock()
	if engineAccount == nil {
		accKb, err := kb.GetByAddress(engineAddress)
		if err != nil {
			panic(err)
		}
		engineAccount = auth.NewBaseAccount(accKb.GetAddress(), nil, accKb.GetPubKey(), 0, 0)
	}
	localSeq := engineAccount.GetSequence()
	lcdGet("auth/accounts/"+engineAddress.String(), &engineAccount)
	// replace seq if sup
	if localSeq > engineAccount.GetSequence() {
		engineAccount.SetSequence(localSeq)
	}
	return engineAccount
}

func createAndSignTx(msgs []sdk.Msg) authtypes.StdTx {
	// retrieve account
	accR := getAccount()
	sequence := accR.GetSequence()
	accR.SetSequence(accR.GetSequence() + 1)

	minGasPrices, err := sdk.ParseDecCoins(cfg.Cosmos.MinGasPrices)
	if err != nil {
		panic(err)
	}

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
	if err != nil {
		panic(err)
	}

	// create StdTx
	stdTx := authtypes.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo)

	// sign StdTx
	signedTx, err := txBuilder.SignStdTx(cfg.Account.Name, cfg.Account.Password, stdTx, false)
	if err != nil {
		panic(err)
	}

	return signedTx
}
