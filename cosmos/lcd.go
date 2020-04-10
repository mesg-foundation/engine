package cosmos

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authExported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// LCD is a simple cosmos LCD client.
type LCD struct {
	endpoint     string
	cdc          *codec.Codec
	kb           keys.Keybase
	chainID      string
	minGasPrices sdktypes.DecCoins
	gasLimit     uint64

	// local state
	accs            map[string]authExported.Account
	accsMutex       sync.Mutex
	getAccountMutex sync.Mutex
}

// NewLCD initializes a cosmos LCD client.
func NewLCD(endpoint string, cdc *codec.Codec, kb keys.Keybase, chainID, minGasPrices string, gasLimit uint64) (*LCD, error) {
	minGasPricesDecoded, err := sdktypes.ParseDecCoins(minGasPrices)
	if err != nil {
		return nil, err
	}
	return &LCD{
		endpoint:     endpoint,
		cdc:          cdc,
		kb:           kb,
		chainID:      chainID,
		minGasPrices: minGasPricesDecoded,
		gasLimit:     gasLimit,
	}, nil
}

// Get executes a get request on the LCD.
func (lcd *LCD) Get(path string, ptr interface{}) error {
	resp, err := http.Get(lcd.endpoint + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request status code is not 2XX, got %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	cosResp := rest.ResponseWithHeight{}
	if err := lcd.cdc.UnmarshalJSON(body, &cosResp); err != nil {
		return err
	}
	if len(cosResp.Result) > 0 {
		if err := lcd.cdc.UnmarshalJSON(cosResp.Result, ptr); err != nil {
			return err
		}
	}
	return nil
}

// Post executes a post request on the LCD.
// The response is expected to be a ResponseWithHeight that contains the expected result.
func (lcd *LCD) Post(path string, req interface{}, ptr interface{}) error {
	cosResp := rest.ResponseWithHeight{}
	if err := lcd.PostBare(path, req, &cosResp); err != nil {
		return err
	}
	if len(cosResp.Result) > 0 {
		if err := lcd.cdc.UnmarshalJSON(cosResp.Result, ptr); err != nil {
			return err
		}
	}
	return nil
}

// PostBare executes a post request on the LCD.
// There is no expectation on the type of response.
func (lcd *LCD) PostBare(path string, req interface{}, ptr interface{}) error {
	reqBody, err := lcd.cdc.MarshalJSON(req)
	if err != nil {
		return err
	}
	resp, err := http.Post(lcd.endpoint+path, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request status code is not 2XX, got %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := lcd.cdc.UnmarshalJSON(body, ptr); err != nil {
		return err
	}
	return nil
}

//BroadcastMsg sign and broadcast a transaction from a message.
func (lcd *LCD) BroadcastMsg(accName, accPassword string, msg sdk.Msg) ([]byte, error) {
	return lcd.BroadcastMsgs(accName, accPassword, []sdk.Msg{msg})
}

//BroadcastMsgs sign and broadcast a transaction from multiple messages.
func (lcd *LCD) BroadcastMsgs(accName, accPassword string, msgs []sdk.Msg) ([]byte, error) {
	tx, err := lcd.createAndSignTx(accName, accPassword, msgs)
	if err != nil {
		return nil, err
	}
	req := authrest.BroadcastReq{
		Tx:   tx,
		Mode: "block", // TODO: should be changed to "sync" and wait for the tx event
	}
	var res sdk.TxResponse
	if err := lcd.PostBare("txs", req, &res); err != nil {
		return nil, err
	}
	if abci.CodeTypeOK != res.Code {
		return nil, fmt.Errorf("transaction returned with invalid code %d: %s", res.Code, res)
	}
	result, err := hex.DecodeString(res.Data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lcd *LCD) getAccount(accName string) (authExported.Account, error) {
	lcd.accsMutex.Lock()
	defer lcd.accsMutex.Unlock()
	if _, ok := lcd.accs[accName]; !ok {
		accKb, err := lcd.kb.Get(accName)
		if err != nil {
			return nil, err
		}
		lcd.accs[accName] = auth.NewBaseAccount(accKb.GetAddress(), nil, accKb.GetPubKey(), 0, 0)
	}
	localSeq := lcd.accs[accName].GetSequence()
	var accR *auth.BaseAccount
	if err := lcd.Get("auth/accounts/"+lcd.accs[accName].GetAddress().String(), &accR); err != nil {
		return nil, err
	}
	lcd.accs[accName] = accR
	// replace seq if sup
	if localSeq > lcd.accs[accName].GetSequence() {
		lcd.accs[accName].SetSequence(localSeq)
	}
	return lcd.accs[accName], nil
}

func (lcd *LCD) createAndSignTx(accName, accPassword string, msgs []sdk.Msg) (authtypes.StdTx, error) {
	// retrieve account
	accR, err := lcd.getAccount(accName)
	if err != nil {
		return authtypes.StdTx{}, err
	}
	sequence := accR.GetSequence()
	accR.SetSequence(accR.GetSequence() + 1)

	// Create TxBuilder
	txBuilder := authtypes.NewTxBuilder(
		authutils.GetTxEncoder(lcd.cdc),
		accR.GetAccountNumber(),
		sequence,
		lcd.gasLimit,
		flags.DefaultGasAdjustment,
		true,
		lcd.chainID,
		"",
		nil,
		lcd.minGasPrices,
	).WithKeybase(lcd.kb)

	// create StdSignMsg
	stdSignMsg, err := txBuilder.BuildSignMsg(msgs)
	if err != nil {
		return authtypes.StdTx{}, err
	}

	// create StdTx
	stdTx := authtypes.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo)

	// sign StdTx
	signedTx, err := txBuilder.SignStdTx(accName, accPassword, stdTx, false)
	if err != nil {
		return authtypes.StdTx{}, err
	}

	return signedTx, nil
}
