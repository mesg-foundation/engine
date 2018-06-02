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
	s := loadService("./service/tests/service-valid")
	assert.NotNil(t, s)
}

func TestTagService(t *testing.T) {
	s := &service.Service{}
	tagService(s, "TestTagService")
	assert.Equal(t, s.Dependencies["service"].Image, "TestTagService")
}

func TestTagServiceWithConfig(t *testing.T) {
	s := &service.Service{
		Configuration: &service.Dependency{
			Command: "xxx",
			Image:   "yyy",
		},
	}
	tagService(s, "TestTagServiceWithConfig")
	assert.Equal(t, s.Dependencies["service"].Image, "TestTagServiceWithConfig")
	assert.Equal(t, s.Dependencies["service"].Command, "xxx")
}

func TestTagServiceWithDependency(t *testing.T) {
	s := &service.Service{
		Dependencies: map[string]*service.Dependency{
			"test": &service.Dependency{
				Image: "xxx",
			},
		},
	}
	tagService(s, "TestTagServiceWithDependency")
	assert.Equal(t, s.Dependencies["service"].Image, "TestTagServiceWithDependency")
	assert.Equal(t, s.Dependencies["test"].Image, "xxx")
}

func TestTagServiceWithDependencyOverride(t *testing.T) {
	s := &service.Service{
		Dependencies: map[string]*service.Dependency{
			"service": &service.Dependency{
				Image: "xxx",
			},
		},
	}
	tagService(s, "TestTagServiceWithDependencyOverride")
	assert.Equal(t, s.Dependencies["service"].Image, "TestTagServiceWithDependencyOverride")
}
