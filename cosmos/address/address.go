package address

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Address interface {
	Equals(sdk.Address) bool
	Empty() bool
	Marshal() ([]byte, error)
	MarshalJSON() ([]byte, error)
	Bytes() []byte
	String() string
	Format(s fmt.State, verb rune)

	// MarshalTo(data []byte) (n int, err error)
	// Unmarshal(data []byte) error
	// Size() int
	// UnmarshalJSON(data []byte) error
}
