package codec

import (
	cosmoscodec "github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
)

// Codec is a general codec where every data structs needs to register to
var Codec = amino.NewCodec()

func init() {
	cosmostypes.RegisterCodec(Codec)
	cosmoscodec.RegisterCrypto(Codec)
}

// RegisterConcrete https://godoc.org/github.com/tendermint/go-amino#Codec.RegisterConcrete
func RegisterConcrete(o interface{}, name string, copts *amino.ConcreteOptions) {
	Codec.RegisterConcrete(o, name, copts)
}

// RegisterInterface https://godoc.org/github.com/tendermint/go-amino#Codec.RegisterInterface
func RegisterInterface(ptr interface{}, iopts *amino.InterfaceOptions) {
	Codec.RegisterInterface(ptr, iopts)
}

// MustMarshalJSON https://godoc.org/github.com/tendermint/go-amino#Codec.MustMarshalJSON
func MustMarshalJSON(o interface{}) []byte {
	return Codec.MustMarshalJSON(o)
}

// MarshalJSON https://godoc.org/github.com/tendermint/go-amino#Codec.MarshalJSON
func MarshalJSON(o interface{}) ([]byte, error) {
	return Codec.MarshalJSON(o)
}

// UnmarshalJSON https://godoc.org/github.com/tendermint/go-amino#Codec.UnmarshalJSON
func UnmarshalJSON(bz []byte, ptr interface{}) error {
	return Codec.UnmarshalJSON(bz, ptr)
}

// UnmarshalBinaryBare https://godoc.org/github.com/tendermint/go-amino#Codec.UnmarshalBinaryBare
func UnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	return Codec.UnmarshalBinaryBare(bz, ptr)
}

// MarshalBinaryBare https://godoc.org/github.com/tendermint/go-amino#Codec.MarshalBinaryBare
func MarshalBinaryBare(o interface{}) ([]byte, error) {
	return Codec.MarshalBinaryBare(o)
}
