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

const bech32PrefixInstAddr = "mesginst"

// Ensure that different address types implement the interface
var _ sdk.Address = InstAddress{}
var _ Address = &InstAddress{}

// InstAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type InstAddress []byte

// InstAddressFromHex creates an InstAddress from a hex string.
func InstAddressFromHex(address string) (addr InstAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding Bech32 address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return InstAddress(bz), nil
}

// InstAddressFromBech32 creates an InstAddress from a Bech32 string.
func InstAddressFromBech32(address string) (addr InstAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return InstAddress{}, nil
	}

	bz, err := sdk.GetFromBech32(address, bech32PrefixInstAddr)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return InstAddress(bz), nil
}

// Returns boolean for whether two InstAddresses are Equal
func (aa InstAddress) Equals(aa2 sdk.Address) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Returns boolean for whether an InstAddress is empty
func (aa InstAddress) Empty() bool {
	if aa == nil {
		return true
	}

	aa2 := InstAddress{}
	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa InstAddress) Marshal() ([]byte, error) {
	return aa, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *InstAddress) Unmarshal(data []byte) error {
	*aa = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa InstAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (aa InstAddress) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *InstAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := InstAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (aa *InstAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := InstAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// Bytes returns the raw address bytes.
func (aa InstAddress) Bytes() []byte {
	return aa
}

// String implements the Stringer interface.
func (aa InstAddress) String() string {
	if aa.Empty() {
		return ""
	}

	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixInstAddr, aa.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa InstAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(aa.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(aa))))
	}
}

// MarshalTo marshal the address in data
func (aa *InstAddress) MarshalTo(data []byte) (int, error) {
	b, err := aa.Marshal()
	if err != nil {
		return 0, err
	}
	return copy(data, b), nil
}

// Size returns the marshaled size
func (aa *InstAddress) Size() int {
	return len(aa.Bytes())
}
