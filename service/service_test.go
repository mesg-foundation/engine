package service

import (
	"errors"
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/assert"
)

var ts = &Service{
	Name: "test-servcie",
	Configuration: &Dependency{
		Key:         mainServiceKey,
		Env:         []string{"foo="},
		Ports:       []string{"80:81", "443"},
		Volumes:     []string{"v1"},
		VolumesFrom: []string{"dummy"},
	},
	Dependencies: []*Dependency{
		{
			Key:         "dummy",
			Volumes:     []string{"v2"},
			VolumesFrom: []string{mainServiceKey},
			Ports:       []string{"8080:8081", "8443"},
		},
	},
	Tasks: []*Task{
		{
			Key: "test",
			Outputs: []*Output{
				{
					Key: "output",
					Data: []*Parameter{
						{
							Key:  "foo",
							Type: paramAnyType,
						},
					},
				},
			},
		},
	},
	Events: []*Event{
		{
			Key: "test",
			Data: []*Parameter{
				{
					Key:  "foo",
					Type: paramAnyType,
				},
			},
		},
	},
}

func TestValidateEventData(t *testing.T) {
	data := map[string]interface{}{"foo": "bar"}
	assert.NoError(t, ts.ValidateEventData("test", data))
	assert.Error(t, ts.ValidateEventData("not-found", data))
}

func TestValidateTaskInputs(t *testing.T) {
	data := map[string]interface{}{"foo": "bar"}
	assert.NoError(t, ts.ValidateTaskInputs("test", data))
	assert.Error(t, ts.ValidateTaskInputs("not-found", data))
}

func TestValidateTaskOutput(t *testing.T) {
	data := map[string]interface{}{"foo": "bar"}
	assert.NoError(t, ts.ValidateTaskOutput("test", "output", data))
	assert.Error(t, ts.ValidateTaskOutput("test", "not-found", data))
	assert.Error(t, ts.ValidateTaskOutput("not-found", "output", data))
}

func TestValidateConfigurationEnv(t *testing.T) {
	assert.NoError(t, ts.validateConfigurationEnv(map[string]string{"foo": "bar"}))
	assert.Error(t, ts.validateConfigurationEnv(map[string]string{"bar": "foo"}))

}

func TestSerivcePorts(t *testing.T) {
	assert.Equal(t, []container.Port{
		{Target: 81, Published: 80},
		{Target: 443, Published: 443},
	}, ts.Configuration.ports())

	assert.Equal(t, []container.Port{
		{Target: 8081, Published: 8080},
		{Target: 8443, Published: 8443},
	}, ts.Dependencies[0].ports())
}

func TestServiceVolumes(t *testing.T) {
	assert.Equal(t, []container.Mount{
		{
			Source: "7631f4220854bd51793cb48df1363bef29013e2cf4bc",
			Target: "v1",
		},
		{
			Source: "7632e460dfba09369509b6946e9924e39e2afe9ddfa1",
			Target: "v2",
		},
	}, ts.volumes(mainServiceKey))

	assert.Equal(t, []container.Mount{
		{
			Source: "7632e460dfba09369509b6946e9924e39e2afe9ddfa1",
			Target: "v2",
		},
		{
			Source: "7631f4220854bd51793cb48df1363bef29013e2cf4bc",
			Target: "v1",
		},
	}, ts.volumes("dummy"))
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
				Type: paramObjectType,
			},
			map[string]interface{}{},
			nil,
		},
		{
			Parameter{
				Type: paramObjectType,
			},
			[]interface{}{},
			errors.New("not an object"),
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.err, tt.param.Validate(tt.arg))
	}
}
