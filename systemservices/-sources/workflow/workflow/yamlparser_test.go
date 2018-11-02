package workflow

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYAMLParser(t *testing.T) {
	yml := `
name: a
description: b
services:
  x: c
  y: d
configs:
  k:
    l: 1
  z: e
when:
  s:
    event:
      t:
        map:
          u: f
          v: true
        execute:
          w: g
`
	def, err := ParseYAML(strings.NewReader(yml))
	require.NoError(t, err)
	require.Equal(t, WorkflowDefinition{
		Name:        "a",
		Description: "b",
		Services: []ServiceDefinition{
			{Name: "x", ID: "c"},
			{Name: "y", ID: "d"},
		},
		Configs: []ConfigDefinition{
			{Key: "k", Value: map[interface{}]interface{}{"l": 1}},
			{Key: "z", Value: "e"},
		},
		Events: []EventDefinition{
			{
				ServiceName: "s",
				EventKey:    "t",
				Map: []MapDefinition{
					{Key: "u", Value: "f"},
					{Key: "v", Value: true},
				},
				Execute: ExecuteDefinition{ServiceName: "w", TaskKey: "g"},
			},
		},
	}, def)
}
