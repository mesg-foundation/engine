package sdk

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/stretchr/testify/require"
)

func TestInjectDefinitionWithConfig(t *testing.T) {
	var (
		command = "xxx"
		s       = &service.Service{}
	)

	injectDefinition(s, &importer.ServiceDefinition{
		Configuration: &importer.Dependency{
			Command: command,
		},
	})
	require.Equal(t, command, s.Configuration.Command)
}

func TestInjectDefinitionWithDependency(t *testing.T) {
	var (
		s     = &service.Service{}
		image = "xxx"
	)

	injectDefinition(s, &importer.ServiceDefinition{
		Dependencies: map[string]*importer.Dependency{
			"test": {
				Image: image,
			},
		},
	})
	require.Equal(t, s.Dependencies[0].Image, image)
}
