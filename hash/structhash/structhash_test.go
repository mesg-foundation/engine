package structhash

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:megacheck
func TestDump(t *testing.T) {
	int1 := int(1)
	tests := []struct {
		v interface{}
		s string
	}{
		{nil, ""},
		{make(chan string), ""},
		{func() {}, ""},
		{false, ""},
		{(*int)(nil), ""},
		{int(0), ""},
		{uint(0), ""},
		{0.0, ""},
		{"", ""},
		{true, "true"},
		{uint32(42), "42"},
		{uint64(42), "42"},
		{int32(-42), "-42"},
		{int64(-42), "-42"},
		{float32(3.14), "3.14"},
		{float64(3.14), "3.14"},
		{interface{}(0), ""},
		{map[int]int(nil), ""},
		{map[int]int{}, ""},
		{map[int]int{0: 0}, ""},
		{map[int]int{0: 0, 1: 1}, "1:1;"},
		{map[int]int{2: 2, 0: 0, 1: 1}, "1:1;2:2;"},
		{map[string]int{"0": 1, "1": 0, "2": 2}, "0:1;2:2;"},
		{map[string]int{"g": 1, "1": 0, "f": 2}, "f:2;g:1;"},
		{[]int(nil), ""},
		{[]*int{nil}, ""},
		{[]int{}, ""},
		{[]int{0}, ""},
		{[]int{0, 1}, "1:1;"},
		{[0]int{}, ""},
		{[1]int{0}, ""},
		{[2]int{0, 1}, "1:1;"},
		{complex(0, 0), ""},
		{(interface{})(nil), ""},
		{(*struct{})(nil), ""},
		{struct{}{}, ""},
		{(chan int)(nil), ""},
		{
			struct {
				a bool           `protobuf:"bytes,1"`
				b int            `protobuf:"bytes,2"`
				c int8           `protobuf:"bytes,3"`
				d int16          `protobuf:"bytes,4"`
				e int32          `protobuf:"bytes,5"`
				f int64          `protobuf:"bytes,6"`
				g uint           `protobuf:"bytes,7"`
				h uint8          `protobuf:"bytes,8"`
				i uint16         `protobuf:"bytes,9"`
				j uint32         `protobuf:"bytes,10"`
				k uint64         `protobuf:"bytes,11"`
				l float32        `protobuf:"bytes,12"`
				m float64        `protobuf:"bytes,13"`
				n []int          `protobuf:"bytes,14"`
				o map[int]int    `protobuf:"bytes,15"`
				p map[string]int `protobuf:"bytes,16"`
				r *int           `protobuf:"bytes,17"`
				s string         `protobuf:"bytes,18"`
			}{},
			"",
		},
		{
			struct {
				a bool           `protobuf:"bytes,1"`
				b int            `protobuf:"bytes,2"`
				c int8           `protobuf:"bytes,3"`
				d int16          `protobuf:"bytes,4"`
				e int32          `protobuf:"bytes,5"`
				f int64          `protobuf:"bytes,6"`
				g uint           `protobuf:"bytes,7"`
				h uint8          `protobuf:"bytes,8"`
				i uint16         `protobuf:"bytes,9"`
				j uint32         `protobuf:"bytes,10"`
				k uint64         `protobuf:"bytes,11"`
				l float32        `protobuf:"bytes,12"`
				m float64        `protobuf:"bytes,13"`
				n []int          `protobuf:"bytes,14"`
				o map[int]int    `protobuf:"bytes,15"`
				p map[string]int `protobuf:"bytes,16"`
				r *int           `protobuf:"bytes,17"`
				s string         `protobuf:"bytes,18"`
			}{
				a: true,
				b: 1,
				c: 1,
				d: 1,
				e: 1,
				f: 1,
				g: 1,
				h: 1,
				i: 1,
				j: 1,
				k: 1,
				l: 1.1,
				m: 1.1,
				n: []int{1, 0, 1},
				o: map[int]int{0: 1, 1: 0},
				p: map[string]int{"0": 1, "1": 0, "2": 2},
				r: &int1,
				s: "1",
			},
			"1:true;2:1;3:1;4:1;5:1;6:1;7:1;8:1;9:1;10:1;11:1;12:1.1;13:1.1;14:CMvDhwpnsTgALRFiAzwi7GR9GUbFwo3xhp9MjifExAW7;15:EUBJnuc9DVgYJtdGABsukyTmTgpoGkUZBEBCZ2GR7HFJ;16:H19PMHCjrY3wgrYpS5qDxGPbfftEpqg68eri5JuPC8qY;17:1;18:1;",
		},
		{
			struct {
				b int `protobuf:"bytes,2"`
				a int `protobuf:"bytes,1"`
			}{2, 1},
			"1:1;2:2;",
		},
		{
			struct {
				a interface{} `protobuf_oneof:"type"`
			}{
				a: struct {
					b int `protobuf:"bytes,1"`
					c int `protobuf:"bytes,2"`
				}{0, 1},
			},
			"BL6qCe785oZSJZ1QQJZuPH9KZE75aSe2wfZUdbyyBAFP",
		},
		{[]string{"foo", "", "bar"}, "0:foo;2:bar;"},
		{
			map[string]string{
				"keyB":  "bar",
				"keyAB": "",
				"keyA":  "foo",
			}, "keyA:foo;keyB:bar;",
		},
		{
			struct {
				a int `protobuf:"bytes,1" hash:"-"`
			}{1},
			"",
		},
	}

	for _, tt := range tests {
		assert.Equalf(t, tt.s, serialize(tt.v), "type %s: %v", reflect.TypeOf(tt.v), tt.v)
	}
}
