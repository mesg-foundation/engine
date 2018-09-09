package service

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mesg-foundation/core/interface/grpc/core"
)

func TestLogReader(t *testing.T) {
	var (
		dep        = "1"
		typ        = core.LogData_Data_Standard
		dataChunk1 = []byte{1}
		dataChunk2 = []byte{2}
	)

	lr := newLogReader(dep, typ)

	go func() {
		lr.process(&core.LogData{
			Data: &core.LogData_Data{
				Dependency: dep,
				Type:       typ,
				Data:       dataChunk1,
			},
		})
		lr.process(&core.LogData{
			Data: &core.LogData_Data{
				Dependency: dep,
				Type:       typ,
				Data:       dataChunk2,
			},
		})
		lr.Close()
	}()

	data, err := ioutil.ReadAll(lr)
	require.NoError(t, err)
	require.Len(t, data, 2)
	require.Equal(t, dataChunk1, []byte{data[0]})
	require.Equal(t, dataChunk2, []byte{data[1]})
}

func TestLogReaderWithWrongDependency(t *testing.T) {
	var (
		dep       = "1"
		dep2      = "2"
		typ       = core.LogData_Data_Standard
		dataChunk = []byte{1}
	)

	lr := newLogReader(dep, typ)

	go func() {
		lr.process(&core.LogData{
			Data: &core.LogData_Data{
				Dependency: dep2,
				Type:       typ,
				Data:       dataChunk,
			},
		})
		lr.Close()
	}()

	data, err := ioutil.ReadAll(lr)
	require.NoError(t, err)
	require.Len(t, data, 0)
}

func TestLogReaderWithWrongType(t *testing.T) {
	var (
		dep       = "1"
		typ       = core.LogData_Data_Standard
		typ2      = core.LogData_Data_Error
		dataChunk = []byte{1}
	)

	lr := newLogReader(dep, typ)

	go func() {
		lr.process(&core.LogData{
			Data: &core.LogData_Data{
				Dependency: dep,
				Type:       typ2,
				Data:       dataChunk,
			},
		})
		lr.Close()
	}()

	data, err := ioutil.ReadAll(lr)
	require.NoError(t, err)
	require.Len(t, data, 0)
}
