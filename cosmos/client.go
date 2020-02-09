package cosmos

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authExported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
	*rpcclient.Local
	kb           keys.Keybase
	chainID      string
	accName      string
	accPassword  string
	minGasPrices sdktypes.DecCoins

	// Local state
	acc             authExported.Account
	getAccountMutex sync.Mutex
	broadcastMutex  sync.Mutex
}

// NewClient returns a rpc tendermint client.
func NewClient(node *node.Node, kb keys.Keybase, chainID, accName, accPassword, minGasPrices string) (*Client, error) {
	minGasPricesDecoded, err := sdktypes.ParseDecCoins(minGasPrices)
	if err != nil {
		return nil, err
	}
	return &Client{
		Local:        rpcclient.NewLocal(node),
		kb:           kb,
		chainID:      chainID,
		accName:      accName,
		accPassword:  accPassword,
		minGasPrices: minGasPricesDecoded,
	}, nil
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
	c.broadcastMutex.Lock() // Lock the whole signature + broadcast of the transaction
	signedTx, err := c.createAndSignTx([]sdktypes.Msg{msg})
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
		return nil, fmt.Errorf("transaction returned with invalid code %d", txres.Code)
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
		return nil, errors.New("i/o timeout")
	}
}

// Stream subscribes to the provided query and returns the hash of the matching ressources.
func (c *Client) Stream(ctx context.Context, query string) (chan hash.Hash, chan error, error) {
	var asciiletters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = asciiletters[rand.Intn(len(asciiletters))]
	}
	subscriber := string(b)
	eventStream, err := c.Subscribe(ctx, subscriber, query, 0)
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
				attrs := event.Events[EventHashType]
				// The following error might be too much as MAYBE if one transaction contains many messages, the events will be merged across the whole transaction
				if len(attrs) != 1 {
					errC <- fmt.Errorf("event %s has %d tag(s), but only 1 is expected", EventHashType, len(attrs))
				}
				for _, attr := range attrs {
					hash, err := hash.Decode(attr)
					if err != nil {
						errC <- err
					} else {
						hashC <- hash
					}
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

// GetAccount returns the local account.
func (c *Client) GetAccount() (authExported.Account, error) {
	c.getAccountMutex.Lock()
	defer c.getAccountMutex.Unlock()
	if c.acc == nil {
		accKb, err := c.kb.Get(c.accName)
		if err != nil {
			return nil, err
		}
		c.acc = auth.NewBaseAccount(
			accKb.GetAddress(),
			nil,
			accKb.GetPubKey(),
			0,
			0,
		)
	}
	localSeq := c.acc.GetSequence()
	accR, err := auth.NewAccountRetriever(c).GetAccount(c.acc.GetAddress())
	if err != nil {
		return nil, err
	}
	c.acc = accR
	// replace seq if sup
	if localSeq > c.acc.GetSequence() {
		c.acc.SetSequence(localSeq)
	}
	return c.acc, nil
}

func (c *Client) createAndSignTx(msgs []sdktypes.Msg) (tenderminttypes.Tx, error) {
	// retrieve account
	accR, err := c.GetAccount()
	if err != nil {
		return nil, err
	}
	sequence := accR.GetSequence()
	accR.SetSequence(accR.GetSequence() + 1)

	// Create TxBuilder
	txBuilder := authtypes.NewTxBuilder(
		authutils.GetTxEncoder(codec.Codec),
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
		_, adjusted, err := authutils.CalculateGas(c.QueryWithData, codec.Codec, txBytes, txBuilder.GasAdjustment())
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
	signedTx, err := txBuilder.SignStdTx(c.accName, c.accPassword, stdTx, false)
	if err != nil {
		return nil, err
	}

	return txBuilder.TxEncoder()(signedTx)
}
