package types

import (
	"testing"

	"github.com/mesg-foundation/engine/hash/structhash"
	"github.com/stretchr/testify/assert"
)

func TestServiceHash(t *testing.T) {
	hashes := make(map[string]bool)
	structs := []*Value{
		{
			Kind: &Value_NullValue{},
		},
		{
			Kind: &Value_NumberValue{},
		},
		{
			Kind: &Value_StringValue{},
		},
		{
			Kind: &Value_BoolValue{},
		},
		{
			Kind: &Value_ListValue{},
		},
	}
	for _, s := range structs {
		hashes[string(structhash.Dump(s))] = true
	}
	assert.Equal(t, len(hashes), len(structs))
}
