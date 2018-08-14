package service

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/xpflag"
	"github.com/stvp/assert"
)

func TestReadJSONFile(t *testing.T) {
	d, e := readJSONFile("")
	assert.Nil(t, e)
	assert.Equal(t, "{}", d)

	d, e = readJSONFile("./doesntexistsfile")
	assert.NotNil(t, e)

	d, e = readJSONFile("./tests/validData.json")
	assert.Nil(t, e)
	assert.Equal(t, "{\"foo\":\"bar\"}", d)
}

func TestTaskKeysFromService(t *testing.T) {
	keys := taskKeysFromService(&service.Service{
		Tasks: map[string]*service.Task{
			"task1": {},
		},
	})
	assert.Equal(t, "task1", keys[0])
}

func TestGetData(t *testing.T) {
	s := service.Service{
		Tasks: map[string]*service.Task{
			"test": {
				Inputs: map[string]*service.Parameter{
					"foo":    {Type: "String"},
					"hello":  {Type: "String"},
					"number": {Type: "Number"},
					"bool":   {Type: "Boolean"},
				},
			},
		},
	}
	var data map[string]string
	cmd := cobra.Command{}
	cmd.Flags().VarP(xpflag.NewStringToStringValue(&data, nil), "data", "", "")
	cmd.Flags().StringP("json", "", "", "")
	cmd.Flags().Set("data", "foo=bar")
	cmd.Flags().Set("data", "hello=world")
	cmd.Flags().Set("data", "number=42")
	cmd.Flags().Set("data", "bool=true")
	res, err := getData(&cmd, "test", &s, data)
	assert.Nil(t, err)
	assert.Equal(t, "{\"bool\":true,\"foo\":\"bar\",\"hello\":\"world\",\"number\":42}", res)
}

func TestGetDataJSON(t *testing.T) {
	s := service.Service{
		Tasks: map[string]*service.Task{
			"test": {
				Inputs: map[string]*service.Parameter{
					"foo": {Type: "String"},
				},
			},
		},
	}
	var data map[string]string
	cmd := cobra.Command{}
	cmd.Flags().VarP(xpflag.NewStringToStringValue(&data, nil), "data", "", "")
	cmd.Flags().StringP("json", "", "", "")
	cmd.Flags().Set("json", "./tests/validData.json")
	res, err := getData(&cmd, "test", &s, data)
	assert.Nil(t, err)
	assert.Equal(t, "{\"foo\":\"bar\"}", res)
}

func TestGetTaskKey(t *testing.T) {
	cmd := cobra.Command{}
	cmd.Flags().StringP("task", "", "", "")
	cmd.Flags().Set("task", "test")
	assert.Equal(t, "test", getTaskKey(&cmd, &service.Service{}))
}
