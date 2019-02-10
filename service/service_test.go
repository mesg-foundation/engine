package service

import (
	"errors"
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/assert"
)

func TestValidateEventData(t *testing.T) {
	s := &Service{
		Events: map[string]*Event{
			"test": {
				Data: map[string]*Parameter{
					"foo": {
						Type: paramAnyType,
					},
				},
			},
		},
	}
	data := map[string]interface{}{"foo": "bar"}
	assert.NoError(t, s.ValidateEventData("test", data))
	assert.Error(t, s.ValidateEventData("not-found", data))
}

func TestValidateTaskInputs(t *testing.T) {
	s := &Service{
		Tasks: map[string]*Task{
			"test": {
				Inputs: map[string]*Parameter{
					"foo": {
						Type: paramAnyType,
					},
				},
			},
		},
	}
	data := map[string]interface{}{"foo": "bar"}
	assert.NoError(t, s.ValidateTaskInputs("test", data))
	assert.Error(t, s.ValidateTaskInputs("not-found", data))
}

func TestValidateTaskOutput(t *testing.T) {
	s := &Service{
		Tasks: map[string]*Task{
			"test": {
				Outputs: map[string]*Output{
					"output": {
						Data: map[string]*Parameter{
							"foo": {
								Type: paramAnyType,
							},
						},
					},
				},
			},
		},
	}
	data := map[string]interface{}{"foo": "bar"}
	assert.NoError(t, s.ValidateTaskOutput("test", "output", data))
	assert.Error(t, s.ValidateTaskOutput("test", "not-found", data))
	assert.Error(t, s.ValidateTaskOutput("not-found", "output", data))
}

func TestValidateConfigurationEnv(t *testing.T) {
	var s Service
	s.Configuration.Env = []string{"foo="}
	assert.NoError(t, s.validateConfigurationEnv(map[string]string{"foo": "bar"}))
	assert.Error(t, s.validateConfigurationEnv(map[string]string{"bar": "foo"}))

}

func TestSerivcePorts(t *testing.T) {
	s := &Service{
		Configuration: Dependency{
			Ports: []string{"80:81", "443"},
		},
		Dependencies: map[string]*Dependency{
			"dummy": {
				Ports: []string{"8080:8081", "8443"},
			},
		},
	}

	assert.Equal(t, []container.Port{
		{Target: 81, Published: 80},
		{Target: 443, Published: 443},
	}, s.ports(MainServiceKey))

	assert.Equal(t, []container.Port{
		{Target: 8081, Published: 8080},
		{Target: 8443, Published: 8443},
	}, s.ports("dummy"))
}

func TestServiceVolumes(t *testing.T) {
	s := &Service{
		Configuration: Dependency{
			Volumes:     []string{"v1"},
			VolumesFrom: []string{"dummy"},
		},
		Dependencies: map[string]*Dependency{
			"dummy": {
				Volumes:     []string{"v2"},
				VolumesFrom: []string{MainServiceKey},
			},
		},
	}

	assert.Equal(t, []container.Mount{
		{
			Source: "76314cf5bc59bee9e1c44c6254b5f84e7f066bd8e5fe",
			Target: "v1",
		},
		{
			Source: "76324cf5bc59bee9e1c44c6254b5f84e7f066bd8e5fe",
			Target: "v2",
		},
	}, s.volumes(MainServiceKey))

	assert.Equal(t, []container.Mount{
		{
			Source: "7632829c3804401b0727f70f73d4415e162400cbe57b",
			Target: "v2",
		},
		{
			Source: "7631829c3804401b0727f70f73d4415e162400cbe57b",
			Target: "v1",
		},
	}, s.volumes("dummy"))
}

func TestVolumeKey(t *testing.T) {
	assert.Equal(t, "766f6c756d65564ca9fb3226c4272a364725c73dbe97ec94a5ee", volumeKey("sid", "dep", "volume"))
}

func TestParametrValidate(t *testing.T) {
	tests := []struct {
		param Parameter
		arg   interface{}
		err   error
	}{
		{
			Parameter{
				Type:     paramStringType,
				Optional: true,
			},
			nil,
			nil,
		},
		{
			Parameter{
				Type: paramStringType,
			},
			nil,
			errors.New("required"),
		},
		{
			Parameter{
				Type: paramStringType,
			},
			"foo",
			nil,
		},
		{
			Parameter{
				Type: paramStringType,
			},
			0,
			errors.New("not a string"),
		},
		{
			Parameter{
				Type: paramNumberType,
			},
			1.0,
			nil,
		},
		{
			Parameter{
				Type: paramNumberType,
			},
			"",
			errors.New("not a number"),
		},
		{
			Parameter{
				Type: paramBooleanType,
			},
			true,
			nil,
		},
		{
			Parameter{
				Type: paramBooleanType,
			},
			"",
			errors.New("not a boolean"),
		},
		{
			Parameter{
				Type:     paramStringType,
				Repeated: true,
			},
			[]interface{}{"foo"},
			nil,
		},
		{
			Parameter{
				Type:     paramStringType,
				Repeated: true,
			},
			"foo",
			errors.New("not an array"),
		},
		{
			Parameter{
				Type: paramOjbectType,
			},
			map[string]interface{}{},
			nil,
		},
		{
			Parameter{
				Type: paramOjbectType,
			},
			[]interface{}{},
			errors.New("not an object"),
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.err, tt.param.Validate(tt.arg))
	}
}
