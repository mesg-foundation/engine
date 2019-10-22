package cosmos

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
)

// GenesisValidator holds the info of a specific validator to use to generate a genesis.
type GenesisValidator struct {
	Name      string
	Address   sdktypes.AccAddress
	Password  string
	Mnemonic  string
	ValPubKey crypto.PubKey
	NodeID    p2p.ID
}

// NewGenesisValidator creates a new validator with an cosmos account, validator and node identity.
func NewGenesisValidator(kb *Keybase, name, password, privValidatorKeyFile, privValidatorStateFile, nodeKeyFile string) (GenesisValidator, error) {
	// create validator account
	mnemonic, err := kb.NewMnemonic()
	if err != nil {
		return GenesisValidator{}, err
	}
	acc, err := kb.CreateAccount(name, mnemonic, "", password, 0, 0)
	if err != nil {
		return GenesisValidator{}, err
	}
	val := privval.LoadOrGenFilePV(privValidatorKeyFile, privValidatorStateFile)
	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyFile)
	if err != nil {
		return GenesisValidator{}, err
	}
	return GenesisValidator{
		Name:      name,
		Password:  password,
		Mnemonic:  mnemonic,
		Address:   acc.GetAddress(),
		ValPubKey: val.GetPubKey(),
		NodeID:    nodeKey.ID(),
	}, nil
}

// Map exports the info of GenesisValidator in a nice way for logrus.
func (validator GenesisValidator) Map() map[string]interface{} {
	return map[string]interface{}{
		"name":     validator.Name,
		"address":  validator.Address,
		"password": validator.Password,
		"mnemonic": validator.Mnemonic,
		"nodeID":   validator.NodeID,
		"peer":     fmt.Sprintf("%s@%s:26656", validator.NodeID, validator.Name),
	}
}

// LoadGenesis loads a genesis from a file.
func LoadGenesis(genesisFile string) (*types.GenesisDoc, error) {
	if _, err := os.Stat(genesisFile); err != nil {
		return nil, err
	}
	return types.GenesisDocFromFile(genesisFile)
}

// GenGenesis generates a new genesis and save it.
func GenGenesis(app *App, kb *Keybase, chainID string, genesisFile string, validators []GenesisValidator) (*types.GenesisDoc, error) {
	msgs := []sdktypes.Msg{}
	for _, validator := range validators {
		// generate msg to add this validator
		msgs = append(msgs, genCreateValidatorMsg(validator.Address, validator.Name, validator.ValPubKey))
	}
	// generate genesis transaction
	b := NewTxBuilder(app.Cdc(), 0, 0, kb, chainID)
	signedMsg, err := b.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}
	validatorTx := authtypes.NewStdTx(signedMsg.Msgs, signedMsg.Fee, []authtypes.StdSignature{}, signedMsg.Memo)
	for _, validator := range validators {
		validatorTx, err = b.SignStdTx(validator.Name, validator.Password, validatorTx, true)
		if err != nil {
			return nil, err
		}
	}

	// generate genesis
	appState, err := genGenesisAppState(app.DefaultGenesis(), app.Cdc(), validatorTx)
	if err != nil {
		return nil, err
	}
	genesis, err := genGenesisDoc(app.Cdc(), appState, chainID, time.Now())
	if err != nil {
		return nil, err
	}
	// save genesis
	if err := genesis.SaveAs(genesisFile); err != nil {
		return nil, err
	}
	return genesis, nil
}

func genGenesisDoc(cdc *codec.Codec, appState map[string]json.RawMessage, chainID string, genesisTime time.Time) (*types.GenesisDoc, error) {
	appStateEncoded, err := cdc.MarshalJSON(appState)
	if err != nil {
		return nil, err
	}
	genesis := &types.GenesisDoc{
		GenesisTime:     genesisTime,
		ChainID:         chainID,
		ConsensusParams: types.DefaultConsensusParams(),
		AppState:        appStateEncoded,
	}
	return genesis, genesis.ValidateAndComplete()
}

func genGenesisAppState(defaultGenesisŚtate map[string]json.RawMessage, cdc *codec.Codec, signedStdTx authtypes.StdTx) (map[string]json.RawMessage, error) {
	genAccs := []genaccounts.GenesisAccount{}

	for _, signer := range signedStdTx.GetSigners() {
		stakes := sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.NewInt(100000000))
		genAcc := genaccounts.NewGenesisAccountRaw(signer, sdktypes.NewCoins(stakes), sdktypes.NewCoins(), 0, 0, "", "")
		if err := genAcc.Validate(); err != nil {
			return nil, err
		}
		genAccs = append(genAccs, genAcc)
	}

	genstate, err := cdc.MarshalJSON(genaccounts.GenesisState(genAccs))
	if err != nil {
		return nil, err
	}
	defaultGenesisŚtate[genaccounts.ModuleName] = genstate

	return genutil.SetGenTxsInAppGenesisState(cdc, defaultGenesisŚtate, []authtypes.StdTx{signedStdTx})
}

func genCreateValidatorMsg(accAddress sdktypes.AccAddress, accName string, valPubKey crypto.PubKey) stakingtypes.MsgCreateValidator {
	return stakingtypes.NewMsgCreateValidator(
		sdktypes.ValAddress(accAddress),
		valPubKey,
		sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.TokensFromConsensusPower(100)),
		stakingtypes.Description{
			Moniker: accName,
			Details: "init-validator",
		},
		stakingtypes.NewCommissionRates(
			sdktypes.ZeroDec(),
			sdktypes.ZeroDec(),
			sdktypes.ZeroDec(),
		),
		sdktypes.NewInt(1),
	)
}
