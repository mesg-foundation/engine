package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/cosmos"
	abci "github.com/tendermint/tendermint/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// Client is a tendermint client with helper functions.
type Client struct {
	rpcclient.Client
	cdc         *codec.Codec
	kb          keys.Keybase
	chainID     string
	address     types.AccAddress
	accName     string
	accPassword string
}

// New returns a rpc tendermint client.
func New(c rpcclient.Client, cdc *codec.Codec, kb keys.Keybase, chainID string, address types.AccAddress, accName, accPassword string) *Client {
	return &Client{
		Client:      c,
		cdc:         cdc,
		kb:          kb,
		chainID:     chainID,
		address:     address,
		accName:     accName,
		accPassword: accPassword,
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
func (c *Client) BuildAndBroadcastMsg(msg sdktypes.Msg, accNumber, accSeq uint64) error {
	txBuilder := cosmos.NewTxBuilder(c.cdc, accNumber, accSeq, c.kb, c.chainID)
	signedTx, err := txBuilder.BuildAndSignStdTx(msg, c.accName, c.accPassword)
	if err != nil {
		return err
	}

	encodedTx, err := txBuilder.Encode(signedTx)
	if err != nil {
		return err
	}

	txres, err := c.BroadcastTxSync(encodedTx)
	if err != nil {
		return err
	}

	if txres.Code != abci.CodeTypeOK {
		return fmt.Errorf("transaction returned with invalid code %d", txres.Code)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	out, err := c.Subscribe(ctx, "", "tx.hash='"+txres.Hash.String()+"'")
	if err != nil {
		return err
	}

	select {
	case <-out:
		return nil
	case <-ctx.Done():
		return errors.New("i/o timeout")
	}
}
