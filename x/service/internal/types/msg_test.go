package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgCreateRequired(t *testing.T) {
	err := MsgCreate{}.ValidateBasic().Error()
	require.Contains(t, err, "Name is a required field")
	require.Contains(t, err, "Source is a required field")
	require.Contains(t, err, "Owner is a required field")
}
