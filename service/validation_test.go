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

func TestValidFile(t *testing.T) {
	_, err := ValidServiceFile("./tests/minimal-valid.yml")
	assert.NotNil(t, err)
}

func TestNonExistingFile(t *testing.T) {
	_, err := ValidServiceFile("./tests/non-existing-file.yml")
	assert.NotNil(t, err)
}

func TestMalFormattedFile(t *testing.T) {
	res, err := ValidServiceFile("./tests/mal-formatted.yml")
	assert.Nil(t, err)
	assert.Equal(t, res.Valid(), false)
	assert.Equal(t, len(res.Errors()), 1)
}

func TestInvalidFile(t *testing.T) {
	_, err := ValidServiceFile("./tests/non-valid.yml")
	assert.NotNil(t, err)
}

func TestInvalidService(t *testing.T) {
	var service *Service
	res, err := service.IsValid()
	assert.Nil(t, err)
	assert.Equal(t, res.Valid(), false)
	assert.Equal(t, len(res.Errors()), 1)
}

func TestMissingName(t *testing.T) {
	service := mergeServices(dependencyValid)
	res, err := service.IsValid()
	assert.Nil(t, err)
	assert.Equal(t, res.Valid(), false)
	assert.Equal(t, len(res.Errors()), 1)
}

func TestMissingDependency(t *testing.T) {
	service := mergeServices(nameValid)
	res, err := service.IsValid()
	assert.Nil(t, err)
	assert.Equal(t, res.Valid(), false)
	assert.Equal(t, len(res.Errors()), 1)
}

func TestInvalidVisibility(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, visibilityInalid)
	res, err := service.IsValid()
	assert.Nil(t, err)
	assert.Equal(t, res.Valid(), false)
	assert.Equal(t, len(res.Errors()), 1)
}

func TestValidVisibility(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, visibilityValid)
	res, _ := service.IsValid()
	assert.Equal(t, res.Valid(), true)
}

func TestInvalidPublish(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, publishInvalid)
	res, err := service.IsValid()
	assert.Nil(t, err)
	assert.Equal(t, res.Valid(), false)
	assert.Equal(t, len(res.Errors()), 1)
}

func TestValidPublish(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, publishValid)
	res, _ := service.IsValid()
	assert.Equal(t, res.Valid(), true)
}

func TestInvalidEvent(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, eventInvalid)
	res, err := service.IsValid()
	assert.Nil(t, err)
	assert.Equal(t, res.Valid(), false)
	assert.True(t, len(res.Errors()) > 0)
}

func TestValidEvent(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, eventValid)
	res, _ := service.IsValid()
	assert.Equal(t, res.Valid(), true)
}

func TestInvalidTask(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, taskInvalid)
	res, err := service.IsValid()
	assert.Nil(t, err)
	assert.Equal(t, res.Valid(), false)
	assert.True(t, len(res.Errors()) > 0)
}

func TestValidTask(t *testing.T) {
	service := mergeServices(nameValid, dependencyValid, taskValid)
	res, _ := service.IsValid()
	assert.Equal(t, res.Valid(), true)
}
