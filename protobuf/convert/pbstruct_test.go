package convert

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertPbStruct(t *testing.T) {
	got := map[string]interface{}{
		"nil":         nil,
		"bool":        true,
		"int":         int(1),
		"int8":        int8(1),
		"int16":       int16(1),
		"int32":       int32(1),
		"int64":       int64(1),
		"uint":        uint(1),
		"uint8":       uint8(1),
		"uint16":      uint16(1),
		"uint32":      uint32(1),
		"uint64":      uint64(1),
		"float32":     float32(1),
		"float64":     float64(1),
		"string":      "string",
		"error":       errors.New("err"),
		"intslice":    []int{1},
		"stringslice": []string{"string"},
		"intarray":    [1]int{1},
		"stringarray": [1]string{"string"},
		"map": map[string]interface{}{
			"nil":         nil,
			"bool":        true,
			"int":         int(1),
			"uint":        uint(1),
			"float64":     float64(1),
			"string":      "string",
			"error":       errors.New("err"),
			"intslice":    []int{1},
			"stringslice": []string{"string"},
		},
		"struct": struct {
			Nil         *bool
			Bool        bool
			Int         int
			Uint        uint
			Float64     float64
			String      string
			Error       error
			IntSlice    []int
			StringSlice []string
		}{
			nil,
			true,
			int(1),
			uint(1),
			float64(1),
			"string",
			errors.New("err"),
			[]int{1},
			[]string{"string"},
		},
	}

	want := map[string]interface{}{
		"nil":         nil,
		"bool":        true,
		"int":         float64(1),
		"int8":        float64(1),
		"int16":       float64(1),
		"int32":       float64(1),
		"int64":       float64(1),
		"uint":        float64(1),
		"uint8":       float64(1),
		"uint16":      float64(1),
		"uint32":      float64(1),
		"uint64":      float64(1),
		"float32":     float64(1),
		"float64":     float64(1),
		"string":      "string",
		"error":       "err",
		"intslice":    []interface{}{float64(1)},
		"stringslice": []interface{}{"string"},
		"intarray":    []interface{}{float64(1)},
		"stringarray": []interface{}{"string"},
		"map": map[string]interface{}{
			"nil":         nil,
			"bool":        interface{}(true),
			"int":         float64(1),
			"uint":        float64(1),
			"float64":     float64(1),
			"string":      "string",
			"error":       "err",
			"intslice":    []interface{}{float64(1)},
			"stringslice": []interface{}{"string"},
		},
		"struct": map[string]interface{}{
			"Nil":         nil,
			"Bool":        true,
			"Int":         float64(1),
			"Uint":        float64(1),
			"Float64":     float64(1),
			"String":      "string",
			"Error":       "err",
			"IntSlice":    []interface{}{float64(1)},
			"StringSlice": []interface{}{"string"},
		},
	}
	require.Equal(t, want, PbStructToMap(MapToPbStruct(got)))
}
