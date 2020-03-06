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

const bech32PrefixOwnAddr = "mesgown"

// Ensure that different address types implement the interface
var _ sdk.Address = OwnAddress{}

// OwnAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type OwnAddress []byte

// OwnAddressFromHex creates an OwnAddress from a hex string.
func OwnAddressFromHex(address string) (addr OwnAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding Bech32 address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return OwnAddress(bz), nil
}

// OwnAddressFromBech32 creates an OwnAddress from a Bech32 string.
func OwnAddressFromBech32(address string) (addr OwnAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return OwnAddress{}, nil
	}

	bz, err := sdk.GetFromBech32(address, bech32PrefixOwnAddr)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return OwnAddress(bz), nil
}

// Returns boolean for whether two OwnAddresses are Equal
func (aa OwnAddress) Equals(aa2 sdk.Address) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Returns boolean for whether an OwnAddress is empty
func (aa OwnAddress) Empty() bool {
	if aa == nil {
		return true
	}

	aa2 := OwnAddress{}
	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa OwnAddress) Marshal() ([]byte, error) {
	return aa, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *OwnAddress) Unmarshal(data []byte) error {
	*aa = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa OwnAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (aa OwnAddress) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *OwnAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := OwnAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (aa *OwnAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := OwnAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// Bytes returns the raw address bytes.
func (aa OwnAddress) Bytes() []byte {
	return aa
}

// String implements the Stringer interface.
func (aa OwnAddress) String() string {
	if aa.Empty() {
		return ""
	}

	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixOwnAddr, aa.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa OwnAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(aa.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(aa))))
	}
}
