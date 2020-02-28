package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	DefaultMinPrice   = "10000atto"
)

var (
	// KeyMinPrice key for the parameter MinPrice
	KeyMinPrice = []byte("MinPrice")
)

// Params - used for initializing default parameter for instance at genesis
type Params struct {
	MinPrice string `json:"min_price" yaml:"min_price"` // min price to pay for an execution
}

// NewParams creates a new Params object
func NewParams(minPrice string) Params {
	return Params{
		MinPrice: minPrice,
	}
}

// ParamKeyTable for auth module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Params:
	MinPrice:			%s
	`, p.MinPrice)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMinPrice, &p.MinPrice, validateMinPrice),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(DefaultMinPrice)
}

func validateMinPrice(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	_, err := sdk.ParseCoins(v)
	return err
}
