package executionsdk

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/require"
)

// TODO: restore test after refactor create of cosmos node for testing.

// func newTesting(t *testing.T) (*Execution, *apiTesting) {
// 	serviceStore, err := store.NewLevelDBStore(servicedbname)
// 	require.NoError(t, err)
// 	db := database.NewServiceDB(serviceStore, codec.New())

// 	instanceStore, err := store.NewLevelDBStore(instdbname)
// 	require.NoError(t, err)
// 	instDB := database.NewInstanceDB(instanceStore, codec.New())
// 	instance := instancesdk.New(nil)

// 	execDB, err := database.NewExecutionDB(execdbname)
// 	require.NoError(t, err)

// 	processDB, err := database.NewProcessDB(processdbname)
// 	require.NoError(t, err)
// 	process := processesdk.New(instance, processDB)

// 	sdk := New(pubsub.New(0), nil, instance, process, execDB)

// 	return sdk, &apiTesting{
// 		T:           t,
// 		serviceDB:   db,
// 		executionDB: execDB,
// 		instanceDB:  instDB,
// 		processDB:   processDB,
// 	}
// }

// var hs1 = hash.Int(1)

// var testService = &service.Service{
// 	Name: "1",
// 	Sid:  "2",
// 	Hash: hs1,
// 	Tasks: []*service.Service_Task{
// 		{Key: "4"},
// 	},
// 	Dependencies: []*service.Service_Dependency{
// 		{Key: "5"},
// 	},
// }

// func TestGet(t *testing.T) {
// 	sdk, at := newTesting(t)
// 	defer at.close()
// 	exec := execution.New(nil, nil, nil, nil, "", "", nil, nil)
// 	require.NoError(t, sdk.execDB.Save(exec))
// 	got, err := sdk.Get(exec.Hash)
// 	require.NoError(t, err)
// 	require.True(t, exec.Equal(got))
// }

// func TestGetStream(t *testing.T) {
// 	sdk, at := newTesting(t)
// 	defer at.close()

// 	exec := execution.New(nil, nil, nil, nil, "", "", nil, nil)
// 	exec.Status = execution.Status_InProgress

// 	require.NoError(t, sdk.execDB.Save(exec))

// 	stream := sdk.GetStream(nil)
// 	defer stream.Close()

// 	go sdk.ps.Pub(exec, streamTopic)
// 	exec.Status = execution.Status_Failed
// 	exec.Error = "exec-error"
// 	require.Equal(t, exec, <-stream.C)
// }

// func TestExecute(t *testing.T) {
// var testInstance = &instance.Instance{
// 	Hash:        hash.Int(2),
// 	ServiceHash: hs1,
// }

// 	sdk, at := newTesting(t)
// 	defer at.close()

// 	require.NoError(t, at.serviceDB.Save(testService))
// 	require.NoError(t, at.instanceDB.Save(testInstance))

// 	_, err := sdk.Execute(nil, testInstance.Hash, hash.Int(1), nil, "", testService.Tasks[0].Key, nil, nil)
// 	require.NoError(t, err)

// 	// not existing instance
// 	_, err = sdk.Execute(nil, hash.Int(3), hash.Int(1), nil, "", testService.Tasks[0].Key, nil, nil)
// 	require.Error(t, err)
// }

// func TestExecuteInvalidTaskKey(t *testing.T) {
// 	sdk, at := newTesting(t)
// 	defer at.close()

// 	require.NoError(t, at.serviceDB.Save(testService))

// 	_, err := sdk.Execute(nil, hs1, hash.Int(1), nil, "", "-", nil, nil)
// 	require.Error(t, err)
// }

func TestSubTopic(t *testing.T) {
	require.Equal(t, subTopic(hash.Hash{0}), "1.Execution")
}
