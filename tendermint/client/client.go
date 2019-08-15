package client

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type Client struct {
	rpcclient.Client
	cdc *codec.Codec
}

func New(c rpcclient.Client, cdc *codec.Codec) *Client {
	return &Client{
		Client: c,
		cdc:    cdc,
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

func (c *Client) SetService(service *service.Service) error {
	return nil
}

func (c *Client) RemoveService(hash hash.Hash) error {
	return nil
}

func (c *Client) GetService(hash hash.Hash) (*service.Service, error) {
	result, err := c.ABCIQuery("custom/serviceapp/service", nil)
	if err != nil {
		return nil, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New(resp.Log)
	}

	var service service.Service
	if err := c.cdc.UnmarshalJSON(resp.Value, &service); err != nil {
		return nil, err
	}
	return &service, nil
}

func (c *Client) ListServices() ([]*service.Service, error) {
	result, err := c.ABCIQuery("custom/serviceapp/services", nil)
	if err != nil {
		return nil, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New(resp.Log)
	}

	var services []*service.Service
	if err := c.cdc.UnmarshalJSON(resp.Value, &services); err != nil {
		return nil, err
	}
	return services, nil
}
