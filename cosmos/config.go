package cosmos

import sdktypes "github.com/cosmos/cosmos-sdk/types"

// See github.com/cosmos/cosmos-sdk/types/address.go
const (
	// Bech32MainPrefix defines the main Bech32 prefix
	Bech32MainPrefix = "mesgtest"

	// CoinType is the mesg registered coin type from https://github.com/satoshilabs/slips/blob/master/slip-0044.md.
	CoinType = 470

	// BIP44Prefix is the parts of the BIP44 HD path that are fixed by
	// what we used during the fundraiser.
	FullFundraiserPath = "44'/470'/0'/0/0"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32MainPrefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32MainPrefix + sdktypes.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixOperator + sdktypes.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixConsensus + sdktypes.PrefixPublic
)

// CustomizeConfig customizes the cosmos application like addresses prefixes and coin type
func CustomizeConfig() {
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.SetFullFundraiserPath(FullFundraiserPath)
	config.SetCoinType(CoinType)
	config.Seal()
}
