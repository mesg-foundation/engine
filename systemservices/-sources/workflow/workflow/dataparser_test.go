package workflow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataParser(t *testing.T) {
	expectations := []struct {
		parser        *dataParser
		value         interface{}
		expectedValue interface{}
		err           error
	}{
		{
			&dataParser{},
			"1",
			"1",
			nil,
		},
		{
			&dataParser{},
			1,
			1,
			nil,
		},
		{
			&dataParser{
				configs: []ConfigDefinition{
					{Key: "a", Value: 1},
				},
			},
			"$configs.a",
			1,
			nil,
		},
		{
			&dataParser{
				configs: []ConfigDefinition{
					{Key: "a", Value: map[string]interface{}{"b": 1}},
				},
			},
			"$configs.a.b",
			1,
			nil,
		},
		{
			&dataParser{
				configs: []ConfigDefinition{
					{Key: "a", Value: 1},
				},
			},
			"$configs.b",
			nil,
			&invalidVarErr{variable: "$configs.b"},
		},
		{
			&dataParser{
				data: map[string]interface{}{
					"a": 1,
				},
			},
			"$data.a",
			1,
			nil,
		},
		{
			&dataParser{
				data: map[string]interface{}{
					"a": map[string]interface{}{"b": 1},
				},
			},
			"$data.a.b",
			1,
			nil,
		},
		{
			&dataParser{
				services: []ServiceDefinition{
					{Name: "x", ID: "a"},
				},
			},
			"$services.x",
			"a",
			nil,
		},
	}

	for _, expectation := range expectations {
		expectedValue, err := expectation.parser.Parse(expectation.value)
		require.Equal(t, expectation.err, err)
		require.Equal(t, expectation.expectedValue, expectedValue)
	}
}
