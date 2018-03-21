package service

import (
	"testing"

	"github.com/imdario/mergo"
	"github.com/stvp/assert"
)

var (
	nameValid       = Service{Name: "ValidName"}
	nameInvalid     = Service{Name: ""}
	dependencyValid = Service{
		Dependencies: Dependencies{
			"x": Dependency{Image: "nginx"},
		},
	}
	dependencyInvalid = Service{
		Dependencies: Dependencies{
			"x": Dependency{Image: ""},
		},
	}
	visibilityValid  = Service{Visibility: vALL}
	visibilityInalid = Service{Visibility: "INVALID"}
	publishValid     = Service{Publish: PublishAll}
	publishInvalid   = Service{Publish: "INVALID"}

	eventValid = Service{
		Events: Events{
			"x": Event{
				Data: Parameters{
					"y": Parameter{
						Type: "String",
					},
				},
			},
		},
	}
	eventInvalid = Service{
		Events: Events{
			"x": Event{
				Data: Parameters{
					"y": Parameter{
						Type: "INVALID_TYPE",
					},
				},
			},
		},
	}
	taskValid = Service{
		Tasks: Tasks{
			"x": Task{
				Outputs: Events{
					"y": Event{
						Data: Parameters{
							"z": Parameter{
								Type: "String",
							},
						},
					},
				},
			},
		},
	}
	taskInvalid = Service{
		Tasks: Tasks{
			"x": Task{},
		},
	}
)

func mergeServices(services ...Service) (service Service) {
	for _, s := range services {
		mergo.Merge(&service, s)
	}
	return
}

func TestInvalidService(t *testing.T) {
	var service *Service
	valid, errs := service.IsValid()
	assert.Equal(t, valid, false)
	assert.Equal(t, len(errs), 1)
}

func TestMissingName(t *testing.T) {
	service := mergeServices(dependencyValid)
	valid, errs := service.IsValid()
	assert.Equal(t, valid, false)
	assert.Equal(t, len(errs), 1)
}

func TestMissingDependency(t *testing.T) {
	service := mergeServices(nameValid)
	valid, errs := service.IsValid()
	assert.Equal(t, valid, false)
	assert.Equal(t, len(errs), 1)
}

func TestInvalidVisibility(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, visibilityInalid)
	valid, errs := service.IsValid()
	assert.Equal(t, valid, false)
	assert.Equal(t, len(errs), 1)
}

func TestValidVisibility(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, visibilityValid)
	valid, _ := service.IsValid()
	assert.Equal(t, valid, true)
}

func TestInvalidPublish(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, publishInvalid)
	valid, errs := service.IsValid()
	assert.Equal(t, valid, false)
	assert.Equal(t, len(errs), 1)
}

func TestValidPublish(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, publishValid)
	valid, _ := service.IsValid()
	assert.Equal(t, valid, true)
}

func TestInvalidEvent(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, eventInvalid)
	valid, errs := service.IsValid()
	assert.Equal(t, valid, false)
	assert.True(t, len(errs) > 0)
}

func TestValidEvent(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, eventValid)
	valid, _ := service.IsValid()
	assert.Equal(t, valid, true)
}

func TestInvalidTask(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, taskInvalid)
	valid, errs := service.IsValid()
	assert.Equal(t, valid, false)
	assert.True(t, len(errs) > 0)
}

func TestValidTask(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, taskValid)
	valid, _ := service.IsValid()
	assert.Equal(t, valid, true)
}
