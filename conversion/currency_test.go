package conversion

import (
	"testing"

	"github.com/stvp/assert"
)

func TestFromStringInvalidFormat(t *testing.T) {
	_, err := FromString("42.3dw2 MESG")
	assert.NotNil(t, err, "The error should not be nil")
}

func TestFromStringInvalidCurrency(t *testing.T) {
	_, err := FromString("42.32 ME3feSG")
	assert.NotNil(t, err, "The error should not be nil")
}

func TestFromStringUnitMESG(t *testing.T) {
	amount, err := FromString("42.32 MESG")
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, amount.Value, uint(42.32*float64(MESG)))
}

func TestFromStringUnitBASE(t *testing.T) {
	amount, err := FromString("32 BASE")
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, amount.Value, uint(32))
}

func TestFromStringNoUnit(t *testing.T) {
	amount, err := FromString("42.32")
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, amount.Value, uint(42.32*float64(MESG)))
}

func TestConvertFromMESGToMESG(t *testing.T) {
	amount, err := FromString("42.32 MESG")
	value := amount.Convert(MESG)
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, value, 42.32)
}

func TestConvertFromMESGToBASE(t *testing.T) {
	amount, err := FromString("42.32 MESG")
	value := amount.Convert(BASE)
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, value, 42.32*float64(MESG))
}

func TestFromStringInvalidNumberWithUnitBASE(t *testing.T) {
	_, err := FromString("42.32 BASE")
	assert.NotNil(t, err, "The error should be nil")
}

func TestConvertFromBASEToBASE(t *testing.T) {
	amount, err := FromString("42 BASE")
	value := amount.Convert(BASE)
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, value, float64(42))
}

func TestConvertFromBASEToMESG(t *testing.T) {
	amount, err := FromString("42 BASE")
	value := amount.Convert(MESG)
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, value, float64(42)/float64(MESG))
}
