package cosmos

import (
	"encoding/json"
	"os"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// GenesisValidator holds the info of a specific validator to use to generate a genesis.
type GenesisValidator struct {
	Name      string
	Password  string
	ValPubKey crypto.PubKey
	NodeID    p2p.ID
}

// NewGenesisValidator creates a new validator with an cosmos account, validator and node identity.
func NewGenesisValidator(kb *Keybase, name, password, privValidatorKeyFile, privValidatorStateFile, nodeKeyFile string) (GenesisValidator, error) {
	val := privval.LoadOrGenFilePV(privValidatorKeyFile, privValidatorStateFile)
	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyFile)
	if err != nil {
		return GenesisValidator{}, err
	}
	return GenesisValidator{
		Name:      name,
		Password:  password,
		ValPubKey: val.GetPubKey(),
		NodeID:    nodeKey.ID(),
	}, nil
}

// GenesisExist returns true if the genesis file already exists.
func GenesisExist(genesisFile string) bool {
	_, err := os.Stat(genesisFile)
	return !os.IsNotExist(err)
}

// LoadGenesis loads a genesis from a file.
func LoadGenesis(genesisFile string) (*tmtypes.GenesisDoc, error) {
	return tmtypes.GenesisDocFromFile(genesisFile)
}

// GenGenesis generates a new genesis and save it.
func GenGenesis(kb *Keybase, defaultGenesisŚtate map[string]json.RawMessage, chainID, initialBalances, validatorDelegationCoin, genesisFile string, validators []GenesisValidator) (*tmtypes.GenesisDoc, error) {
	valDelCoin, err := sdktypes.ParseCoin(validatorDelegationCoin)
	if err != nil {
		return nil, err
	}
	msgs := []sdktypes.Msg{}
	for _, validator := range validators {
		// get account
		acc, err := kb.Get(validator.Name)
		if err != nil {
			return nil, err
		}
		// generate msg to add this validator
		msgs = append(msgs, genCreateValidatorMsg(acc.GetAddress(), validator.Name, valDelCoin, validator.ValPubKey))
	}
	// generate genesis transaction
	accNumber := uint64(0)
	sequence := uint64(0)
	b := NewTxBuilder(accNumber, sequence, kb, chainID, sdktypes.DecCoins{})
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
	appState, err := genGenesisAppState(defaultGenesisŚtate, validatorTx, initialBalances)
	if err != nil {
		return nil, err
	}
	genesis, err := genGenesisDoc(appState, chainID, time.Now())
	if err != nil {
		return nil, err
	}
	// save genesis
	if err := genesis.SaveAs(genesisFile); err != nil {
		return nil, err
	}
	return genesis, nil
}

func genGenesisDoc(appState map[string]json.RawMessage, chainID string, genesisTime time.Time) (*tmtypes.GenesisDoc, error) {
	appStateEncoded, err := codec.MarshalJSON(appState)
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

func genGenesisAppState(defaultGenesisŚtate map[string]json.RawMessage, signedStdTx authtypes.StdTx, initialBalances string) (map[string]json.RawMessage, error) {
	genAccs := []genaccounts.GenesisAccount{}
	for _, signer := range signedStdTx.GetSigners() {
		initialB, err := sdktypes.ParseCoins(initialBalances)
		if err != nil {
			return nil, err
		}
		genAcc := genaccounts.NewGenesisAccountRaw(signer, initialB, sdktypes.NewCoins(), 0, 0, "", "")
		if err := genAcc.Validate(); err != nil {
			return nil, err
		}
		genAccs = append(genAccs, genAcc)
	}
	genstate, err := codec.MarshalJSON(genaccounts.GenesisState(genAccs))
	if err != nil {
		return nil, err
	}
	defaultGenesisŚtate[genaccounts.ModuleName] = genstate
	return genutil.SetGenTxsInAppGenesisState(codec.Codec, defaultGenesisŚtate, []authtypes.StdTx{signedStdTx})
}

func genCreateValidatorMsg(accAddress sdktypes.AccAddress, accName string, validatorDelegationCoin sdktypes.Coin, valPubKey crypto.PubKey) stakingtypes.MsgCreateValidator {
	return stakingtypes.NewMsgCreateValidator(
		sdktypes.ValAddress(accAddress),
		valPubKey,
		validatorDelegationCoin,
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
