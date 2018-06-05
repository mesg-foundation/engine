package cmdService

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestDefaultPath(t *testing.T) {
	assert.Equal(t, defaultPath([]string{}), "./")
	assert.Equal(t, defaultPath([]string{"foo"}), "foo")
	assert.Equal(t, defaultPath([]string{"foo", "bar"}), "foo")
}

func TestLoadService(t *testing.T) {
	s, path := loadService("./service/tests/service-valid")
	assert.NotNil(t, s)
	assert.Equal(t, path, "./service/tests/service-valid")
}

func TestInjectConfigurationInDependencies(t *testing.T) {
	s := &service.Service{}
	injectConfigurationInDependencies(s, "TestInjectConfigurationInDependencies")
	assert.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependencies")
}

func TestInjectConfigurationInDependenciesWithConfig(t *testing.T) {
	s := &service.Service{
		Configuration: &service.Dependency{
			Command: "xxx",
			Image:   "yyy",
		},
	}
	injectConfigurationInDependencies(s, "TestInjectConfigurationInDependenciesWithConfig")
	assert.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependenciesWithConfig")
	assert.Equal(t, s.Dependencies["service"].Command, "xxx")
}

func TestInjectConfigurationInDependenciesWithDependency(t *testing.T) {
	s := &service.Service{
		Dependencies: map[string]*service.Dependency{
			"test": &service.Dependency{
				Image: "xxx",
			},
		},
	}
	injectConfigurationInDependencies(s, "TestInjectConfigurationInDependenciesWithDependency")
	assert.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependenciesWithDependency")
	assert.Equal(t, s.Dependencies["test"].Image, "xxx")
}

func TestInjectConfigurationInDependenciesWithDependencyOverride(t *testing.T) {
	s := &service.Service{
		Dependencies: map[string]*service.Dependency{
			"service": &service.Dependency{
				Image: "xxx",
			},
		},
	}
	injectConfigurationInDependencies(s, "TestInjectConfigurationInDependenciesWithDependencyOverride")
	assert.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependenciesWithDependencyOverride")
}
