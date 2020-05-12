package cosmos

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	// CosmosBech32MainPrefix defines the main Bech32 prefix.
	CosmosBech32MainPrefix = "mesg"

	// CosmosCoinType is the mesg registered coin type from https://github.com/satoshilabs/slips/blob/master/slip-0044.md.
	CosmosCoinType = uint32(470)

	// FullFundraiserPath is the parts of the BIP44 HD path that are fixed by what we used during the fundraiser.
	FullFundraiserPath = "44'/470'/0'/0/0"

	// GasAdjustment is a multiplier to make sure transactions have enough gas when gas are estimated.
	GasAdjustment = 1.5

	// DefaultAlgo for create account.
	DefaultAlgo = keys.Secp256k1
)

// InitConfig sets the bech32 prefix and HDPath to cosmos config.
func InitConfig() {
	// See github.com/cosmos/cosmos-sdk/types/address.go
	const (
		// bech32PrefixAccAddr defines the Bech32 prefix of an account's address
		bech32PrefixAccAddr = CosmosBech32MainPrefix
		// bech32PrefixAccPub defines the Bech32 prefix of an account's public key
		bech32PrefixAccPub = CosmosBech32MainPrefix + sdktypes.PrefixPublic
		// bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
		bech32PrefixValAddr = CosmosBech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixOperator
		// bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
		bech32PrefixValPub = CosmosBech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixOperator + sdktypes.PrefixPublic
		// bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
		bech32PrefixConsAddr = CosmosBech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixConsensus
		// bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
		bech32PrefixConsPub = CosmosBech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixConsensus + sdktypes.PrefixPublic
	)

	cfg := sdktypes.GetConfig()
	cfg.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)
	cfg.SetFullFundraiserPath(FullFundraiserPath)
	cfg.SetCoinType(CosmosCoinType)
	cfg.Seal()
}
