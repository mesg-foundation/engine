package types

import (
	"testing"

	"github.com/mesg-foundation/engine/hash/structhash"
	"github.com/stretchr/testify/assert"
)

func TestServiceHash(t *testing.T) {
	hashes := make(map[string]bool)

	structs := []Struct{
		{
			Fields: map[string]*Value{},
		},
		{
			Fields: map[string]*Value{
				"v": {Kind: &Value_NullValue{}},
			},
		},
		{
			Fields: map[string]*Value{
				"v": {Kind: &Value_NumberValue{}},
			},
		},
		{
			Fields: map[string]*Value{
				"v": {Kind: &Value_StringValue{}},
			},
		},
		{
			Fields: map[string]*Value{
				"v": {Kind: &Value_BoolValue{}},
			},
		},
		{
			Fields: map[string]*Value{
				"v": {Kind: &Value_StructValue{}},
			},
		},
		{
			Fields: map[string]*Value{
				"v": {Kind: &Value_ListValue{}},
			},
		},
	}

	for _, s := range structs {
		hashes[string(structhash.Dump(s))] = true
	}

	assert.Equal(t, len(hashes), len(structs))
}
