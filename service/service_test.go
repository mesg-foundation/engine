package service

import (
	"testing"

	"github.com/mesg-foundation/core/service/importer"
	"github.com/stretchr/testify/require"
)

func TestErrNotDefinedEnv(t *testing.T) {
	require.Equal(t, ErrNotDefinedEnv{[]string{"A", "B"}}.Error(),
		`environment variable(s) "A, B" not defined in mesg.yml (under configuration.env key)`)
}

func TestInjectDefinitionWithConfig(t *testing.T) {
	var (
		command = "xxx"
		s       = &Service{}
	)

	s.InjectDefinition(&importer.ServiceDefinition{
		Configuration: &importer.Dependency{
			Command: command,
		},
	})
	require.Equal(t, command, s.Configuration().Command)
}

func TestInjectDefinitionWithDependency(t *testing.T) {
	var (
		s     = &Service{}
		image = "xxx"
	)

	s.InjectDefinition(&importer.ServiceDefinition{
		Dependencies: map[string]*importer.Dependency{
			"test": {
				Image: image,
			},
		},
	})
	require.Equal(t, s.Dependencies[0].Image, image)
}
