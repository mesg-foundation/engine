package cosmos

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/x/xreflect"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/node"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tenderminttypes "github.com/tendermint/tendermint/types"
)

// Client is a tendermint client with helper functions.
type Client struct {
	rpcclient.Client
	kb      keys.Keybase
	chainID string
}

// NewClient returns a rpc tendermint client.
func NewClient(node *node.Node, kb keys.Keybase, chainID string) *Client {
	return &Client{
		Client:  rpcclient.NewLocal(node),
		kb:      kb,
		chainID: chainID,
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
func (c *Client) BuildAndBroadcastMsg(msg sdktypes.Msg, accName, accPassword string) (*abci.ResponseDeliverTx, error) {
	info, err := c.kb.Get(accName)
	if err != nil {
		return nil, err
	}

	accRetriever := auth.NewAccountRetriever(c)
	accNumber, accSeq := uint64(0), uint64(0)
	err = accRetriever.EnsureExists(info.GetAddress())
	if err == nil {
		accNumber, accSeq, err = accRetriever.GetAccountNumberSequence(info.GetAddress())
		if err != nil {
			return nil, err
		}
	}

	txBuilder := NewTxBuilder(accNumber, accSeq, c.kb, c.chainID)

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

// Stream subscribes to the provided query and returns the hash of the matching ressources.
func (c *Client) Stream(query string) (<-chan hash.Hash, func() error, error) {
	ctx, cancel := context.WithCancel(context.Background())
	// TODO: randStringRunes seems to complicated for this random string
	subscriber := randStringRunes(8) // generates a new subscriber each time to be able to subscribe to the same query multiple time
	eventStream, err := c.Subscribe(ctx, subscriber, query)
	if err != nil {
		return nil, func() error {
			cancel()
			return nil
		}, err
	}
	closer := func() error {
		err := c.Unsubscribe(ctx, subscriber, query)
		if err != nil {
			return err
		}
		cancel()
		return nil
	}
	hashChan := make(chan hash.Hash)
	go func() {
		defer close(hashChan)
		for event := range eventStream {
			ressHashes := event.Events[eventHashKey()]
			if len(ressHashes) != 1 {
				// TODO: remove panic
				panic(errors.New("there is no or more than one " + eventHashKey()))
			}
			hash, err := hash.Decode(ressHashes[0])
			if err != nil {
				// TODO: remove panic
				panic(err)
			}
			hashChan <- hash
		}
	}()
	return hashChan, closer, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	nbr := len(letterRunes)
	for i := range b {
		b[i] = letterRunes[rand.Intn(nbr)]
	}
	return string(b)
}
