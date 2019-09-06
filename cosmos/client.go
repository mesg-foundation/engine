package cosmos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
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

// New returns a rpc tendermint client.
func NewClient(node *node.Node, cdc *codec.Codec, kb keys.Keybase, chainID string) *Client {
	return &Client{
		Client:  rpcclient.NewLocal(node),
		cdc:     cdc,
		kb:      kb,
		chainID: chainID,
	}
}

// QueryWithData is abci.query wraper with errors check.
// It retruns slice of bytes, height and an error.
func (c *Client) QueryWithData(path string, data []byte) ([]byte, int64, error) {
	result, err := c.ABCIQuery(path, data)
	if err != nil {
		return nil, 0, err
	}
	resp := result.Response
	if !resp.IsOK() {
		return nil, 0, errors.New(resp.Log)
	}
	return resp.Value, resp.Height, nil
}

// BuildAndBroadcastMsg builds and signs message and broadcast it to node.
func (c *Client) BuildAndBroadcastMsg(msg sdktypes.Msg, accName, accPassword string, accNumber, accSeq uint64) (*tenderminttypes.TxResult, error) {
	txBuilder := NewTxBuilder(c.cdc, accNumber, accSeq, c.kb, c.chainID)
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

	out, err := c.Subscribe(ctx, "", "tx.hash='"+txres.Hash.String()+"'")
	if err != nil {
		return nil, err
	}

	select {
	case result := <-out:
		data, ok := result.Data.(tenderminttypes.EventDataTx)
		if !ok {
			return nil, errors.New("result data is not the right type")
		}
		return &data.TxResult, nil
	case <-ctx.Done():
		return nil, errors.New("i/o timeout")
	}
}
