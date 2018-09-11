package service

import (
	"testing"

	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/x/xpflag"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestReadJSONFile(t *testing.T) {
	d, e := readJSONFile("")
	require.Nil(t, e)
	require.Equal(t, "{}", d)

	d, e = readJSONFile("./doesntexistsfile")
	require.NotNil(t, e)

	d, e = readJSONFile("./tests/validData.json")
	require.Nil(t, e)
	require.Equal(t, "{\"foo\":\"bar\"}", d)
}

func TestTaskKeysFromService(t *testing.T) {
	keys := taskKeysFromService(&core.Service{
		Tasks: map[string]*core.Task{
			"task1": {},
		},
	})
	require.Equal(t, "task1", keys[0])
}

func TestGetData(t *testing.T) {
	s := core.Service{
		Tasks: map[string]*core.Task{
			"test": {
				Inputs: map[string]*core.Parameter{
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
	require.Nil(t, err)
	require.Equal(t, "{\"bool\":true,\"foo\":\"bar\",\"hello\":\"world\",\"number\":42}", res)
}

func TestGetDataError(t *testing.T) {
	s := core.Service{
		Tasks: map[string]*core.Task{
			"test": {
				Inputs: map[string]*core.Parameter{
					"bool": {Type: "Boolean"},
				},
			},
		},
	}
	var data map[string]string
	cmd := cobra.Command{}
	cmd.Flags().VarP(xpflag.NewStringToStringValue(&data, nil), "data", "", "")
	cmd.Flags().StringP("json", "", "", "")
	cmd.Flags().Set("data", "bool=hello")
	_, err := getData(&cmd, "test", &s, data)
	require.NotNil(t, err)
}

func TestGetDataJSON(t *testing.T) {
	s := core.Service{
		Tasks: map[string]*core.Task{
			"test": {
				Inputs: map[string]*core.Parameter{
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
	require.Nil(t, err)
	require.Equal(t, "{\"foo\":\"bar\"}", res)
}

func TestGetTaskKey(t *testing.T) {
	cmd := cobra.Command{}
	cmd.Flags().StringP("task", "", "", "")
	cmd.Flags().Set("task", "test")
	require.Equal(t, "test", getTaskKey(&cmd, &core.Service{}))
}

func TestExecutePreRun(t *testing.T) {
	cmd := cobra.Command{}
	cmd.Flags().StringP("data", "", "", "")
	cmd.Flags().StringP("json", "", "", "")
	cmd.Flags().Set("json", "test")
	executePreRun(&cmd, []string{})
}
