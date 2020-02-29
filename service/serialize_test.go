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
	require.Equal(t, "1:world;5:0:6:0:3:Number;4:true;8:key;;;7:0:3:Number;4:true;8:key;;;8:task;;;6:0:3:0:3:Number;4:true;8:key;;;4:event;;;12:hello;", s.Serialize())
	require.Equal(t, "5Hwubvgm5eDFJyXXpEEnTcYNFR7Jppmq94vmQb7oxXfX", hash.Dump(s).String())
}

func BenchmarkSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.Serialize()
	}
}

// with base58
// BenchmarkSerialize-4   	   45284	     26965 ns/op	    7928 B/op	     216 allocs/op
// BenchmarkSerialize-4   	   44287	     27114 ns/op	    7928 B/op	     216 allocs/op
// BenchmarkSerialize-4   	   45910	     26025 ns/op	    7928 B/op	     216 allocs/op

// without base58
// BenchmarkSerialize-4   	   98306	     11453 ns/op	    5848 B/op	     188 allocs/op
// BenchmarkSerialize-4   	   95521	     11667 ns/op	    5848 B/op	     188 allocs/op
// BenchmarkSerialize-4   	  106831	     11040 ns/op	    5848 B/op	     188 allocs/op

// last one + eliminate empty value as soon as possible
// BenchmarkSerialize-4   	  130498	      8952 ns/op	    4352 B/op	     149 allocs/op
// BenchmarkSerialize-4   	  121597	     10585 ns/op	    4352 B/op	     149 allocs/op
// BenchmarkSerialize-4   	  127125	      9357 ns/op	    4352 B/op	     149 allocs/op

// last one + no sort
// BenchmarkSerialize-4   	  136657	      8135 ns/op	    4000 B/op	     138 allocs/op
// BenchmarkSerialize-4   	  133406	      8202 ns/op	    4000 B/op	     138 allocs/op
// BenchmarkSerialize-4   	  137749	      8333 ns/op	    4000 B/op	     138 allocs/op

// last one + bytes
// BenchmarkSerialize-4   	  234852	      4824 ns/op	    3432 B/op	      86 allocs/op
// BenchmarkSerialize-4   	  224851	      4972 ns/op	    3432 B/op	      86 allocs/op
// BenchmarkSerialize-4   	  219336	      5380 ns/op	    3432 B/op	      86 allocs/op

// last one + buffer pool
// BenchmarkSerialize-4   	  244681	      4571 ns/op	    2328 B/op	      73 allocs/op
// BenchmarkSerialize-4   	  246154	      4622 ns/op	    2328 B/op	      73 allocs/op
// BenchmarkSerialize-4   	  237271	      4530 ns/op	    2328 B/op	      73 allocs/op

// last one + return early on empty pairs
// BenchmarkSerialize-4   	  263187	      4401 ns/op	    2328 B/op	      73 allocs/op
// BenchmarkSerialize-4   	  236986	      4268 ns/op	    2328 B/op	      73 allocs/op
// BenchmarkSerialize-4   	  279604	      4284 ns/op	    2328 B/op	      73 allocs/op

// last one + no sort
// BenchmarkSerialize-4   	  291133	      3688 ns/op	    1976 B/op	      62 allocs/op

// last one + no pairs, direct write to a dedicated buffer (no buffer pool)
// BenchmarkSerialize-4   	  468823	      2347 ns/op	    2000 B/op	      30 allocs/op
// BenchmarkSerialize-4   	  459656	      2329 ns/op	    2000 B/op	      30 allocs/op
// BenchmarkSerialize-4   	  458227	      2260 ns/op	    2000 B/op	      30 allocs/op
