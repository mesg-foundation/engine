package cosmos

import (
	"math/big"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/mesg-foundation/engine/config"
)

func init() {
	// See github.com/cosmos/cosmos-sdk/types/address.go
	const (
		// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
		Bech32PrefixAccAddr = config.CosmosBech32MainPrefix
		// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
		Bech32PrefixAccPub = config.CosmosBech32MainPrefix + sdktypes.PrefixPublic
		// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
		Bech32PrefixValAddr = config.CosmosBech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixOperator
		// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
		Bech32PrefixValPub = config.CosmosBech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixOperator + sdktypes.PrefixPublic
		// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
		Bech32PrefixConsAddr = config.CosmosBech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixConsensus
		// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
		Bech32PrefixConsPub = config.CosmosBech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixConsensus + sdktypes.PrefixPublic
	)

	cfg := sdktypes.GetConfig()
	cfg.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	cfg.SetFullFundraiserPath(config.FullFundraiserPath)
	cfg.SetCoinType(config.CosmosCoinType)
	cfg.Seal()
}

// CustomizeConfig customizes the cosmos application like addresses prefixes and coin type
func CustomizeConfig(engineCfg *config.Config) {
	// From github.com/cosmos/cosmos-sdk/types/staking.go
	// Set the power reduction from token to voting power to 10^18 (number of decimal of the token).
	sdktypes.PowerReduction = sdktypes.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(engineCfg.Cosmos.PowerReduction), nil))

	// From github.com/cosmos/cosmos-sdk/x/staking/types/params.go
	// Override default genesis state of staking module to set the bond denom
	staking.DefaultGenesisState = func() staking.GenesisState {
		state := stakingtypes.DefaultGenesisState()
		state.Params.BondDenom = engineCfg.Cosmos.StakeTokenDenom
		return state
	}
}
