package codec

import "github.com/tendermint/go-amino"

// Codec is a general codec where every data structs needs to register to
var Codec = amino.NewCodec()

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
	if len(bz) == 0 {
		return nil
	}
	return Codec.UnmarshalJSON(bz, ptr)
}
