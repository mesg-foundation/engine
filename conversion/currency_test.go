package conversion

import (
	"testing"

	"github.com/stvp/assert"
)

func TestAmountFromStringInvalidFormat(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42.3dw2 MESG")
	assert.NotNil(t, err, "The error should not be nil")
}

func TestAmountFromStringInvalidCurrency(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42.32 ME3feSG")
	assert.NotNil(t, err, "The error should not be nil")
}

func TestAmountFromStringUnitMESG(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42.32 MESG")
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, amount.Value, uint(42.32*float64(MESG)))
}

func TestAmountFromStringUnitBASE(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("32 BASE")
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, amount.Value, uint(32))
}

func TestAmountFromStringNoUnit(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42.32")
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, amount.Value, uint(42.32*float64(MESG)))
}

func TestAmountConvertFromMESGToMESG(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42.32 MESG")
	value := amount.Convert(MESG)
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, value, 42.32)
}

func TestAmountConvertFromMESGToBASE(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42.32 MESG")
	value := amount.Convert(BASE)
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, value, 42.32*float64(MESG))
}

func TestAmountFromStringInvalidNumberWithUnitBASE(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42.32 BASE")
	assert.NotNil(t, err, "The error should be nil")
}

func TestAmountConvertFromBASEToBASE(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42 BASE")
	value := amount.Convert(BASE)
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, value, float64(42))
}

func TestAmountConvertFromBASEToMESG(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42 BASE")
	value := amount.Convert(MESG)
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, value, float64(42)/float64(MESG))
}

func TestAmountFormatWithBase(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42.44 MESG")
	assert.Nil(t, err, "The error should be nil")
	formattedString := amount.Format(BASE)
	assert.Equal(t, formattedString, "42440000000 BASE")
}

func TestAmountFormatWithMESG(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42 BASE")
	assert.Nil(t, err, "The error should be nil")
	formattedString := amount.Format(MESG)
	assert.Equal(t, formattedString, "0.000000042 MESG")
}

func TestAmountString(t *testing.T) {
	amount := &Amount{}
	err := amount.FromString("42 BASE")
	assert.Nil(t, err, "The error should be nil")
	formattedString := amount.String()
	assert.Equal(t, formattedString, "0.000000042 MESG")
}
