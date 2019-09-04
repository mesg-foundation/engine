package structhash

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha1(t *testing.T) {
	assert.Equal(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", fmt.Sprintf("%x", Sha1(nil)))
}

func TestSha3(t *testing.T) {
	assert.Equal(t, "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26", fmt.Sprintf("%x", Sha3(nil)))
}

func TestMd5(t *testing.T) {
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", fmt.Sprintf("%x", Md5(nil)))
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
		{
			struct {
				a chan int `hash:"name:a"`
			}{},
			"{}",
		},
		{
			struct {
				a int `hash:"omitempty"`
			}{1},
			"{}",
		},
		{
			struct {
				a interface{} `hash:"name:a"`
			}{nil},
			"{a:nil}",
		},
		{
			struct {
				A interface{} `hash:"name:A"`
			}{0},
			"{A:0}",
		},
		{
			struct {
				a int `hash:"name:a"`
			}{0},
			"{a:0}",
		},
		{
			struct {
				a int `hash:"name:a"`
				b int `hash:"name:b"`
			}{0, 1},
			"{a:0,b:1}",
		},
		{
			struct {
				a struct {
					a bool `hash:"name:a"`
				} `hash:"name:a"`
			}{a: struct {
				a bool `hash:"name:a"`
			}{a: false}},
			"{a:{a:false}}",
		},
		{
			struct {
				a *struct {
					b bool `hash:"name:b"`
				} `hash:"name:a"`
			}{a: &struct {
				b bool `hash:"name:b"`
			}{b: false}},
			"{a:{b:false}}",
		},
	}

	for _, tt := range tests {
		s := Dump(tt.v)
		assert.Equalf(t, []byte(tt.s), s, "type %s: %v", reflect.TypeOf(tt.v), tt.v)
	}
}

//nolint:megacheck
func TestTag(t *testing.T) {
	tests := []struct {
		v interface{}
		s string
	}{
		{
			struct {
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
			}{},
			"{}",
		},
		{
			struct {
				a int         `hash:"name:a,omitempty"`
				b uint        `hash:"name:b,omitempty"`
				c bool        `hash:"name:c,omitempty"`
				d string      `hash:"name:d,omitempty"`
				e []int       `hash:"name:e,omitempty"`
				f float64     `hash:"name:f,omitempty"`
				g complex128  `hash:"name:g,omitempty"`
				h interface{} `hash:"name:h,omitempty"`
				i *struct{}   `hash:"name:i,omitempty"`
				j *[]uint     `hash:"name:j,omitempty"`
				k chan int    `hash:"name:k,omitempty"`
			}{},
			"{}",
		},
		{
			struct {
				a int `hash:"name:b"`
			}{},
			"{b:0}",
		},
		{
			struct {
				a int `hash:"name:b"`
				b int `hash:"name:a"`
			}{0, 1},
			"{a:1,b:0}",
		},
		{
			struct {
				a int `hash:"name:b,omitempty"`
				b int `hash:"name:a,omitempty"`
			}{0, 1},
			"{a:1}",
		},
		{
			struct {
				a int
				b int `hash:"-"`
				c int `hash:"name:c"`
			}{1, 2, 3},
			"{c:3}",
		},
	}

	for _, tt := range tests {
		assert.Equalf(t, []byte(tt.s), serialize(tt.v), "type %s: %v", reflect.TypeOf(tt.v), tt.v)
	}
}

func TestTagPanicInvalidOption(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("serialize did not panic with invalid option")
		}
	}()

	serialize(struct {
		a int `hash:"name:a,omitempty,invalid"`
	}{0})
}
