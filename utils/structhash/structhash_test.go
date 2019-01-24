package structhash

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSha1(t *testing.T) {
	require.Equal(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", Sha1Str(nil))
}

func TestMd5(t *testing.T) {
	require.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", Md5Str(nil))
}

//nolint:megacheck
func TestDump(t *testing.T) {
	tests := []struct {
		v interface{}
		s string
	}{
		{nil, ""},
		{make(chan string), ""},
		{func() {}, ""},
		{false, "false"},
		{(*int)(nil), "nil"},
		{int(0), "0"},
		{uint(0), "0"},
		{0.0, "0e+00"},
		{"", `""`},
		{interface{}(0), "0"},
		{map[int]int(nil), "()nil"},
		{map[int]int{}, "()"},
		{map[int]int{0: 0}, "(0:0)"},
		{map[int]int{0: 0, 1: 1}, "(0:0,1:1)"},
		{[]int(nil), "[]nil"},
		{[]*int{nil}, "[nil]"},
		{[]int{}, "[]"},
		{[]int{0}, "[0]"},
		{[]int{0, 1}, "[0,1]"},
		{[0]int{}, "[]"},
		{[1]int{0}, "[0]"},
		{[2]int{0, 1}, "[0,1]"},
		{complex(0, 0), "0e+000e+00i"},
		{(*struct{})(nil), "nil"},
		{struct{}{}, "{}"},
		{struct{ a chan int }{}, "{}"},
		{struct{ a interface{} }{nil}, "{a:nil}"},
		{struct{ A interface{} }{0}, "{A:0}"},
		{struct{ a int }{0}, "{a:0}"},
		{struct{ a, b int }{0, 1}, "{a:0,b:1}"},
		{struct{ a struct{ a bool } }{a: struct{ a bool }{a: false}}, "{a:{a:false}}"},
		{struct{ a *struct{ b bool } }{a: &struct{ b bool }{b: false}}, "{a:{b:false}}"},
	}

	for _, tt := range tests {
		if s := Dump(tt.v); s != tt.s {
			require.Equalf(t, tt.s, s, "type %s: %v", reflect.TypeOf(tt.v), tt.v)
		}
	}
}

//nolint:megacheck
func TestTag(t *testing.T) {
	tests := []struct {
		v interface{}
		s string
	}{
		{struct {
			a int         `hash:"-"`
			b uint        `hash:"-"`
			c bool        `hash:"-"`
			d string      `hash:"-"`
			e []int       `hash:"-"`
			f float64     `hash:"-"`
			g complex128  `hash:"-"`
			h interface{} `hash:"-"`
			i *struct{}   `hash:"-"`
			j *[]uint     `hash:"-"`
			k chan int    `hash:"-"`
		}{}, "{}"},
		{struct {
			a int         `hash:"omitempty"`
			b uint        `hash:"omitempty"`
			c bool        `hash:"omitempty"`
			d string      `hash:"omitempty"`
			e []int       `hash:"omitempty"`
			f float64     `hash:"omitempty"`
			g complex128  `hash:"omitempty"`
			h interface{} `hash:"omitempty"`
			i *struct{}   `hash:"omitempty"`
			j *[]uint     `hash:"omitempty"`
			k chan int    `hash:"omitempty"`
		}{}, "{}"},
		{struct {
			a int `hash:"name:b"`
		}{}, "{b:0}"},
		{struct {
			a int `hash:"name:b"`
			b int `hash:"name:a"`
		}{0, 1}, "{a:1,b:0}"},
		{struct {
			a int `hash:"name:b,omitempty"`
			b int `hash:"name:a,omitempty"`
		}{0, 1}, "{a:1}"},
	}

	for _, tt := range tests {
		if s := serialize(tt.v); string(s) != tt.s {
			require.Equalf(t, tt.s, string(s), "type %s: %v", reflect.TypeOf(tt.v), tt.v)
		}
	}
}
