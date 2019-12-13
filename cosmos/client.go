package cosmos

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/x/xreflect"
	"github.com/mesg-foundation/engine/x/xstrings"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/node"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tenderminttypes "github.com/tendermint/tendermint/types"
)

// Client is a tendermint client with helper functions.
type Client struct {
	rpcclient.Client
	kb          keys.Keybase
	chainID     string
	accName     string
	accPassword string
}

// NewClient returns a rpc tendermint client.
func NewClient(node *node.Node, kb keys.Keybase, chainID, accName, accPassword string) *Client {
	return &Client{
		Client:      rpcclient.NewLocal(node),
		kb:          kb,
		chainID:     chainID,
		accName:     accName,
		accPassword: accPassword,
	}
}

// Query is abci.query wrapper with errors check and decode data.
func (c *Client) Query(path string, qdata, ptr interface{}) error {
	var data []byte
	if !xreflect.IsNil(qdata) {
		b, err := codec.MarshalBinaryBare(qdata)
		if err != nil {
			return err
		}
		data = b
	}
	result, _, err := c.QueryWithData(path, data)
	if err != nil {
		return err
	}
	return codec.UnmarshalBinaryBare(result, ptr)
}

// QueryWithData performs a query to a Tendermint node with the provided path
// and a data payload. It returns the result and height of the query upon success
// or an error if the query fails.
func (c *Client) QueryWithData(path string, data []byte) ([]byte, int64, error) {
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
func (c *Client) BuildAndBroadcastMsg(msg sdktypes.Msg) (*abci.ResponseDeliverTx, error) {
	signedTx, err := c.Sign(msg)
	if err != nil {
		return nil, err
	}

	txres, err := c.BroadcastTxSync(signedTx)
	if err != nil {
		return nil, err
	}

	if txres.Code != abci.CodeTypeOK {
		return nil, fmt.Errorf("transaction returned with invalid code %d", txres.Code)
	}

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
		return nil, errors.New("i/o timeout")
	}
}

// Stream subscribes to the provided query and returns the hash of the matching ressources.
func (c *Client) Stream(ctx context.Context, query string) (chan hash.Hash, chan error, error) {
	subscriber := xstrings.RandASCIILetters(8)
	eventStream, err := c.Subscribe(context.Background(), subscriber, query)
	if err != nil {
		return nil, nil, err
	}

	hashC := make(chan hash.Hash)
	errC := make(chan error)
	go func() {
	loop:
		for {
			select {
			case event := <-eventStream:
				tags := event.Events[EventHashType]
				if len(tags) != 1 {
					errC <- fmt.Errorf("event %s has %d tag(s), but only 1 is expected", EventHashType, len(tags))
					break
				}

				hash, err := hash.Decode(tags[0])
				if err != nil {
					errC <- err
				} else {
					hashC <- hash
				}
			case <-ctx.Done():
				break loop
			}
		}
		close(errC)
		close(hashC)
		c.Unsubscribe(context.Background(), subscriber, query)
	}()
	return hashC, errC, nil
}

var (
	acc             auth.Account
	getAccountMutex sync.Mutex
	accSetSeqMutex  sync.Mutex
)

// GetAccount returns the local account.
func (c *Client) GetAccount() (auth.Account, error) {
	getAccountMutex.Lock()
	defer getAccountMutex.Unlock()
	if acc == nil {
		accKb, err := c.kb.Get(c.accName)
		if err != nil {
			return nil, err
		}
		acc = auth.NewBaseAccount(
			accKb.GetAddress(),
			nil,
			accKb.GetPubKey(),
			AccNumber,
			0,
		)
	}
	localSeq := acc.GetSequence()
	if accR, err := auth.NewAccountRetriever(c).GetAccount(acc.GetAddress()); err == nil {
		acc = accR
		// replace seq if sup
		if localSeq > acc.GetSequence() {
			acc.SetSequence(localSeq)
		}
	}
	return acc, nil
}

// Sign signs a msg and return a tendermint tx.
func (c *Client) Sign(msg sdktypes.Msg) (tenderminttypes.Tx, error) {
	acc, err := c.GetAccount()
	if err != nil {
		return nil, err
	}

	accSetSeqMutex.Lock()
	sequence := acc.GetSequence()
	acc.SetSequence(acc.GetSequence() + 1)
	accSetSeqMutex.Unlock()

	txBuilder := NewTxBuilder(sequence, c.kb, c.chainID)
	signedTx, err := txBuilder.BuildAndSignStdTx(msg, c.accName, c.accPassword)
	if err != nil {
		return nil, err
	}
	return txBuilder.Encode(signedTx)
}
