package cosmos

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/config"
)

// InitConfig sets the bech32 prefix and HDPath to cosmos config.
func InitConfig() {
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

// GasAdjustment is a multiplier to make sure transactions have enough gas when gas are estimated.
const GasAdjustment = 1.5
