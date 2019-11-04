package cosmos

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/node"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tenderminttypes "github.com/tendermint/tendermint/types"
)

// Client is a tendermint client with helper functions.
type Client struct {
	rpcclient.Client
	cdc     *codec.Codec
	kb      keys.Keybase
	chainID string
}

// AccSeq is a sequence number that could be used to build the transaction.
// It should be increased after every usage.
var AccSeq uint64

// NewClient returns a rpc tendermint client.
func NewClient(node *node.Node, cdc *codec.Codec, kb keys.Keybase, chainID string) *Client {
	return &Client{
		Client:  rpcclient.NewLocal(node),
		cdc:     cdc,
		kb:      kb,
		chainID: chainID,
	}
}

// Query is abci.query wrapper with errors check and decode data.
func (c *Client) Query(path string, data cmn.HexBytes, ptr interface{}) error {
	result, err := c.ABCIQuery(path, data)
	if err != nil {
		return err
	}
	if !result.Response.IsOK() {
		return errors.New(result.Response.Log)
	}
	return c.cdc.UnmarshalJSON(result.Response.Value, ptr)
}

// BuildAndBroadcastMsg builds and signs message and broadcast it to node.
func (c *Client) BuildAndBroadcastMsg(msg sdktypes.Msg, accName, accPassword string, accNumber uint64) (*abci.ResponseDeliverTx, error) {
	txBuilder := NewTxBuilder(c.cdc, accNumber, atomic.AddUint64(&AccSeq, 1), c.kb, c.chainID)
	// TODO: cannot sign 2 tx at the same time. Maybe keybase cannot be access at the same time. Add a lock?
	signedTx, err := txBuilder.BuildAndSignStdTx(msg, accName, accPassword)
	if err != nil {
		return nil, err
	}

	encodedTx, err := txBuilder.Encode(signedTx)
	if err != nil {
		return nil, err
	}

	txres, err := c.BroadcastTxSync(encodedTx)
	if err != nil {
		return nil, err
	}

	if txres.Code != abci.CodeTypeOK {
		return nil, fmt.Errorf("transaction returned with invalid code %d", txres.Code)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	subscriber := "engine"
	query := tenderminttypes.EventQueryTxFor(encodedTx).String()
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
		return nil, errors.New("i/o timeout")
	}
}
