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
)

var (
	// KeyMinters key for the parameter Minters
	KeyMinters = []byte("Minters")

	// DefaultMinters is the default value of Minters
	DefaultMinters = []sdk.AccAddress{}
)

// Params - used for initializing default parameter for instance at genesis
type Params struct {
	Minters []sdk.AccAddress `json:"minters" yaml:"minters"`
}

// NewParams creates a new Params object
func NewParams(minters []sdk.AccAddress) Params {
	return Params{
		Minters: minters,
	}
}

// ParamKeyTable for auth module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Params:	
	Minters:			%s	
	`, p.Minters)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMinters, &p.Minters, validateMinters),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(DefaultMinters)
}

func validateMinters(i interface{}) error {
	minters, ok := i.([]sdk.AccAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, minter := range minters {
		if err := sdk.VerifyAddressFormat(minter); err != nil {
			return err
		}
	}
	return nil
}
