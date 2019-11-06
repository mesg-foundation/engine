package executionsdk

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	processesdk "github.com/mesg-foundation/engine/sdk/process"
	"github.com/stretchr/testify/require"
)

const (
	servicedbname = "service.db.test"
	instdbname    = "instance.db.test"
	execdbname    = "exec.db.test"
	processdbname = "process.db.test"
)

type apiTesting struct {
	*testing.T
	serviceDB   *database.ServiceDB
	executionDB *database.LevelDBExecutionDB
	instanceDB  *database.InstanceDB
	processDB   *database.LevelDBProcessDB
}

func (t *apiTesting) close() {
	require.NoError(t, t.serviceDB.Close())
	require.NoError(t, t.executionDB.Close())
	require.NoError(t, t.instanceDB.Close())
	require.NoError(t, t.processDB.Close())
	require.NoError(t, os.RemoveAll(servicedbname))
	require.NoError(t, os.RemoveAll(execdbname))
	require.NoError(t, os.RemoveAll(instdbname))
	require.NoError(t, os.RemoveAll(processdbname))
}

func newTesting(t *testing.T) (*Execution, *apiTesting) {
	serviceStore, err := store.NewLevelDBStore(servicedbname)
	require.NoError(t, err)
	db := database.NewServiceDB(serviceStore, codec.New())

	instanceStore, err := store.NewLevelDBStore(instdbname)
	require.NoError(t, err)
	instDB := database.NewInstanceDB(instanceStore, codec.New())
	instance := instancesdk.New(nil)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	processDB, err := database.NewProcessDB(processdbname)
	require.NoError(t, err)
	process := processesdk.New(instance, processDB)

	sdk := New(pubsub.New(0), nil, instance, process, execDB)

	return sdk, &apiTesting{
		T:           t,
		serviceDB:   db,
		executionDB: execDB,
		instanceDB:  instDB,
		processDB:   processDB,
	}
}

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

func TestGet(t *testing.T) {
	sdk, at := newTesting(t)
	defer at.close()
	exec := execution.New(nil, nil, nil, nil, "", "", nil, nil)
	require.NoError(t, sdk.execDB.Save(exec))
	got, err := sdk.Get(exec.Hash)
	require.NoError(t, err)
	require.True(t, exec.Equal(got))
}

func TestGetStream(t *testing.T) {
	sdk, at := newTesting(t)
	defer at.close()

	exec := execution.New(nil, nil, nil, nil, "", "", nil, nil)
	exec.Status = execution.Status_InProgress

	require.NoError(t, sdk.execDB.Save(exec))

	stream := sdk.GetStream(nil)
	defer stream.Close()

	go sdk.ps.Pub(exec, streamTopic)
	exec.Status = execution.Status_Failed
	exec.Error = "exec-error"
	require.Equal(t, exec, <-stream.C)
}

// TODO: restore test after refactor create of cosmos node for testing.
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
