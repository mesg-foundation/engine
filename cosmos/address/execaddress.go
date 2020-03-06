package address

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/bech32"
	"gopkg.in/yaml.v2"
)

const bech32PrefixExecAddr = "mesgexec"

// Ensure that different address types implement the interface
var _ sdk.Address = ExecAddress{}

// ExecAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type ExecAddress []byte

// ExecAddressFromHex creates an ExecAddress from a hex string.
func ExecAddressFromHex(address string) (addr ExecAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding Bech32 address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return ExecAddress(bz), nil
}

// ExecAddressFromBech32 creates an ExecAddress from a Bech32 string.
func ExecAddressFromBech32(address string) (addr ExecAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return ExecAddress{}, nil
	}

	bz, err := sdk.GetFromBech32(address, bech32PrefixExecAddr)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return ExecAddress(bz), nil
}

// Returns boolean for whether two ExecAddresses are Equal
func (aa ExecAddress) Equals(aa2 sdk.Address) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Returns boolean for whether an ExecAddress is empty
func (aa ExecAddress) Empty() bool {
	if aa == nil {
		return true
	}

	aa2 := ExecAddress{}
	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa ExecAddress) Marshal() ([]byte, error) {
	return aa, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *ExecAddress) Unmarshal(data []byte) error {
	*aa = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa ExecAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (aa ExecAddress) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *ExecAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := ExecAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (aa *ExecAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := ExecAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// Bytes returns the raw address bytes.
func (aa ExecAddress) Bytes() []byte {
	return aa
}

// String implements the Stringer interface.
func (aa ExecAddress) String() string {
	if aa.Empty() {
		return ""
	}

	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixExecAddr, aa.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa ExecAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(aa.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(aa))))
	}
}
