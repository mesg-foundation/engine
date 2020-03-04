package service

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/require"
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

func TestSerialize(t *testing.T) {
	require.Equal(t, "1:world;5:0:6:0:3:Number;4:true;8:key;;;7:0:3:Number;4:true;8:key;;;8:task;;;6:0:3:0:3:Number;4:true;8:key;;;4:event;;;12:hello;", s.HashSerialize())
	require.Equal(t, "5Hwubvgm5eDFJyXXpEEnTcYNFR7Jppmq94vmQb7oxXfX", hash.Dump(s).String())
}

func BenchmarkSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.HashSerialize()
	}
}
