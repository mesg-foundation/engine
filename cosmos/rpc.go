package cosmos

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authExported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/mesg-foundation/engine/ext/xreflect"
	abci "github.com/tendermint/tendermint/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tenderminttypes "github.com/tendermint/tendermint/types"
)

// RPC is a tendermint rpc client with helper functions.
type RPC struct {
	rpcclient.Client
	cdc          *codec.Codec
	kb           keys.Keybase
	chainID      string
	minGasPrices sdktypes.DecCoins

	// local state
	accs           map[string]authExported.Account
	accsMutex      sync.Mutex
	broadcastMutex sync.Mutex
}

// NewRPC returns a rpc tendermint client.
func NewRPC(client rpcclient.Client, cdc *codec.Codec, kb keys.Keybase, chainID, minGasPrices string) (*RPC, error) {
	minGasPricesDecoded, err := sdktypes.ParseDecCoins(minGasPrices)
	if err != nil {
		return nil, err
	}
	return &RPC{
		Client:       client,
		cdc:          cdc,
		kb:           kb,
		chainID:      chainID,
		minGasPrices: minGasPricesDecoded,
	}, nil
}

// Codec returns the codec used by RPC.
func (c *RPC) Codec() *codec.Codec {
	return c.cdc
}

// QueryJSON is abci.query wrapper with errors check and decode data.
func (c *RPC) QueryJSON(path string, qdata, ptr interface{}) error {
	var data []byte
	if !xreflect.IsNil(qdata) {
		b, err := c.cdc.MarshalJSON(qdata)
		if err != nil {
			return err
		}
		data = b
	}

	result, _, err := c.QueryWithData(path, data)
	if err != nil {
		return err
	}
	return c.cdc.UnmarshalJSON(result, ptr)
}

// QueryWithData performs a query to a Tendermint node with the provided path
// and a data payload. It returns the result and height of the query upon success
// or an error if the query fails.
func (c *RPC) QueryWithData(path string, data []byte) ([]byte, int64, error) {
	result, err := c.ABCIQuery(path, data)
	if err != nil {
		return nil, 0, err
	}
	resp := result.Response
	if !resp.IsOK() {
		return nil, resp.Height, errors.New(resp.Log)
	}
	return resp.Value, resp.Height, nil
}

// BuildAndBroadcastMsg builds and signs message and broadcast it to node.
func (c *RPC) BuildAndBroadcastMsg(accName, accPassword string, msg sdktypes.Msg) (*abci.ResponseDeliverTx, error) {
	c.broadcastMutex.Lock() // Lock the whole signature + broadcast of the transaction
	signedTx, err := c.CreateAndSignTx(accName, accPassword, []sdktypes.Msg{msg})
	if err != nil {
		c.broadcastMutex.Unlock()
		return nil, err
	}

	txres, err := c.BroadcastTxSync(signedTx)
	c.broadcastMutex.Unlock()
	if err != nil {
		return nil, err
	}

	if txres.Code != abci.CodeTypeOK {
		return nil, fmt.Errorf("transaction returned with invalid code %d: %s", txres.Code, txres.Log)
	}

	// TODO: 20*time.Second should not be hardcoded here
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	subscriber := "engine"
	query := tenderminttypes.EventQueryTxFor(signedTx).String()
	out, err := c.Subscribe(ctx, subscriber, query)
	if err != nil {
		return nil, err
	}
	defer c.Unsubscribe(ctx, subscriber, query)

	select {
	case result := <-out:
		data, ok := result.Data.(tenderminttypes.EventDataTx)
		if !ok {
			return nil, errors.New("result data is not the right type")
		}
		if data.TxResult.Result.IsErr() {
			return nil, fmt.Errorf("an error occurred in transaction: %s", data.TxResult.Result.Log)
		}
		return &data.TxResult.Result, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("reach timeout for listening for transaction result: %w", ctx.Err())
	}
}

// GetAccount returns an account.
func (c *RPC) GetAccount(accName string) (authExported.Account, error) {
	c.accsMutex.Lock()
	defer c.accsMutex.Unlock()
	if _, ok := c.accs[accName]; !ok {
		accKb, err := c.kb.Get(accName)
		if err != nil {
			return nil, err
		}
		c.accs[accName] = auth.NewBaseAccount(accKb.GetAddress(), nil, accKb.GetPubKey(), 0, 0)
	}
	localSeq := c.accs[accName].GetSequence()
	accR, err := auth.NewAccountRetriever(c).GetAccount(c.accs[accName].GetAddress())
	if err != nil {
		return nil, err
	}
	c.accs[accName] = accR
	// replace seq if sup
	if localSeq > c.accs[accName].GetSequence() {
		c.accs[accName].SetSequence(localSeq)
	}
	return c.accs[accName], nil
}

// CreateAndSignTx build and sign a msg with client account.
func (c *RPC) CreateAndSignTx(accName, accPassword string, msgs []sdktypes.Msg) (tenderminttypes.Tx, error) {
	// retrieve account
	accR, err := c.GetAccount(accName)
	if err != nil {
		return nil, err
	}
	sequence := accR.GetSequence()
	accR.SetSequence(accR.GetSequence() + 1)

	// Create TxBuilder
	txBuilder := authtypes.NewTxBuilder(
		authutils.GetTxEncoder(c.cdc),
		accR.GetAccountNumber(),
		sequence,
		flags.DefaultGasLimit,
		flags.DefaultGasAdjustment,
		true,
		c.chainID,
		"",
		nil,
		c.minGasPrices,
	).WithKeybase(c.kb)

	// calculate gas
	if txBuilder.SimulateAndExecute() {
		txBytes, err := txBuilder.BuildTxForSim(msgs)
		if err != nil {
			return nil, err
		}
		_, adjusted, err := authutils.CalculateGas(c.QueryWithData, c.cdc, txBytes, txBuilder.GasAdjustment())
		if err != nil {
			return nil, err
		}
		txBuilder = txBuilder.WithGas(adjusted)
	}

	// create StdSignMsg
	stdSignMsg, err := txBuilder.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}

	// create StdTx
	stdTx := authtypes.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo)

	// sign StdTx
	signedTx, err := txBuilder.SignStdTx(accName, accPassword, stdTx, false)
	if err != nil {
		return nil, err
	}

	return txBuilder.TxEncoder()(signedTx)
}
