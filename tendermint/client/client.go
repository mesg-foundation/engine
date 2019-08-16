package client

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/mesg-foundation/engine/hash"
	pbtypes "github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/tendermint/app/serviceapp"
	"github.com/mesg-foundation/engine/tendermint/txbuilder"
	"github.com/sirupsen/logrus"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type account struct {
	number uint64
	seq    uint64
}

type Client struct {
	rpcclient.Client
	cdc         *codec.Codec
	kb          keys.Keybase
	chainID     string
	address     types.AccAddress
	accName     string
	accPassword string

	accounts map[string]account
}

func New(c rpcclient.Client, cdc *codec.Codec, kb keys.Keybase, chainID string, address types.AccAddress, accName, accPassword string) *Client {
	return &Client{
		Client:      c,
		cdc:         cdc,
		kb:          kb,
		chainID:     chainID,
		address:     address,
		accName:     accName,
		accPassword: accPassword,
		accounts:    make(map[string]account),
	}
}

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

func (c *Client) SetService(service *pbtypes.Service) error {
	msg := serviceapp.NewMsgSetService(service, c.address)

	acc, ok := c.accounts[c.address.String()]
	if !ok {
		number, seq, err := authtypes.NewAccountRetriever(c).GetAccountNumberSequence(c.address)
		if err != nil {
			return err
		}
		acc.number, acc.seq = number, seq
	}
	logrus.Warning("SEQUENCE NUMBER >>>>>>>>> ", acc.seq, " | ", c.chainID)

	txBuilder := txbuilder.NewTxBuilder(c.cdc, acc.number, acc.seq, c.kb, c.chainID)
	signedTx, err := txBuilder.Create(msg, c.accName, c.accPassword)
	if err != nil {
		return err
	}

	encodedTx, err := txBuilder.Encode(signedTx)
	if err != nil {
		return err
	}

	if _, err := c.BroadcastTxSync(encodedTx); err != nil {
		return err
	}

	acc.seq++
	c.accounts[c.address.String()] = acc
	return nil
}

func (c *Client) RemoveService(hash hash.Hash) error {
	msg := serviceapp.NewMsgRemoveService(hash, c.address)

	acc, ok := c.accounts[c.address.String()]
	if !ok {
		number, seq, err := authtypes.NewAccountRetriever(c).GetAccountNumberSequence(c.address)
		if err != nil {
			return err
		}
		acc.number, acc.seq = number, seq
	}

	txBuilder := txbuilder.NewTxBuilder(c.cdc, acc.number, acc.seq, c.kb, c.chainID)
	signedTx, err := txBuilder.Create(msg, c.accName, c.accPassword)
	if err != nil {
		return err
	}

	encodedTx, err := txBuilder.Encode(signedTx)
	if err != nil {
		return err
	}

	if _, err := c.BroadcastTxSync(encodedTx); err != nil {
		return err
	}

	acc.seq++
	c.accounts[c.address.String()] = acc
	return nil
}

func (c *Client) GetService(hash hash.Hash) (*pbtypes.Service, error) {
	result, err := c.ABCIQuery(fmt.Sprintf("custom/serviceapp/service/%s", hash), nil)
	if err != nil {
		return nil, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New(resp.Log)
	}

	var service pbtypes.Service
	if err := c.cdc.UnmarshalJSON(resp.Value, &service); err != nil {
		return nil, err
	}
	return &service, nil
}

func (c *Client) ListServices() ([]*pbtypes.Service, error) {
	result, err := c.ABCIQuery("custom/serviceapp/services", nil)
	if err != nil {
		return nil, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New(resp.Log)
	}

	var services []*pbtypes.Service
	if err := c.cdc.UnmarshalJSON(resp.Value, &services); err != nil {
		return nil, err
	}
	return services, nil
}
