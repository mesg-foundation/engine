package service

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

var sp = []*Service_Parameter{
	{
		Key:      "key",
		Type:     "Number",
		Optional: true,
	},
}

var s = &Service{
	Sid:  "hello",
	Name: "world",
	Events: []*Service_Event{
		{
			Key:  "event",
			Data: sp,
		},
	},
	Tasks: []*Service_Task{
		{
			Key:     "task",
			Inputs:  sp,
			Outputs: sp,
		},
	},
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "1:world;5:0:6:0:3:Number;4:true;8:key;;;7:0:3:Number;4:true;8:key;;;8:task;;;6:0:3:0:3:Number;4:true;8:key;;;4:event;;;12:hello;", s.HashSerialize())
	require.Equal(t, "cosmos18lrpl337qh9t8ceeuta48k5sz3pjetqs9zjkqk", sdk.AccAddress(crypto.AddressHash([]byte(s.HashSerialize()))).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.HashSerialize()
	}
}