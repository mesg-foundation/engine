package executionsdk

import (
	"os"
	"testing"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/container/mocks"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	workflowsdk "github.com/mesg-foundation/engine/sdk/workflow"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/store"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	servicedbname  = "service.db.test"
	instdbname     = "instance.db.test"
	execdbname     = "exec.db.test"
	workflowdbname = "workflow.db.test"
)

type apiTesting struct {
	*testing.T
	serviceDB   *database.ServiceDB
	executionDB *database.LevelDBExecutionDB
	instanceDB  *database.LevelDBInstanceDB
	workflowDB  *database.LevelDBWorkflowDB
}

func (t *apiTesting) close() {
	require.NoError(t, t.serviceDB.Close())
	require.NoError(t, t.executionDB.Close())
	require.NoError(t, t.instanceDB.Close())
	require.NoError(t, t.workflowDB.Close())
	require.NoError(t, os.RemoveAll(servicedbname))
	require.NoError(t, os.RemoveAll(execdbname))
	require.NoError(t, os.RemoveAll(instdbname))
	require.NoError(t, os.RemoveAll(workflowdbname))
}

func newTesting(t *testing.T) (*Execution, *apiTesting) {
	container := &mocks.Container{}
	s, err := leveldb.OpenFile(servicedbname, nil)
	require.NoError(t, err)
	db := database.NewServiceDB(store.NewLevelDBStore(s))
	service := servicesdk.NewDeprecated(container, db)

	instDB, err := database.NewInstanceDB(instdbname)
	require.NoError(t, err)
	instance := instancesdk.New(container, service, instDB, "", "")

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	workflowDB, err := database.NewWorkflowDB(workflowdbname)
	require.NoError(t, err)
	workflow := workflowsdk.New(instance, workflowDB)

	sdk := New(pubsub.New(0), service, instance, workflow, execDB)

	return sdk, &apiTesting{
		T:           t,
		serviceDB:   db,
		executionDB: execDB,
		instanceDB:  instDB,
		workflowDB:  workflowDB,
	}
}

var testService = &service.Service{
	Name: "1",
	Sid:  "2",
	Hash: hash.Int(1),
	Tasks: []*service.Task{
		{Key: "4"},
	},
	Dependencies: []*service.Dependency{
		{Key: "5"},
	},
}

var testInstance = &instance.Instance{
	Hash:        hash.Int(2),
	ServiceHash: testService.Hash,
}

func TestGet(t *testing.T) {
	sdk, at := newTesting(t)
	defer at.close()
	exec := execution.New(nil, nil, nil, nil, "", "", nil, nil)
	require.NoError(t, sdk.execDB.Save(exec))
	got, err := sdk.Get(exec.Hash)
	require.NoError(t, err)
	require.Equal(t, exec, got)
}

func TestGetStream(t *testing.T) {
	sdk, at := newTesting(t)
	defer at.close()

	exec := execution.New(nil, nil, nil, nil, "", "", nil, nil)
	exec.Status = execution.InProgress

	require.NoError(t, sdk.execDB.Save(exec))

	stream := sdk.GetStream(nil)
	defer stream.Close()

	go sdk.ps.Pub(exec, streamTopic)
	exec.Status = execution.Failed
	exec.Error = "exec-error"
	require.Equal(t, exec, <-stream.C)
}

func TestExecute(t *testing.T) {
	sdk, at := newTesting(t)
	defer at.close()

	require.NoError(t, at.serviceDB.Save(testService))
	require.NoError(t, at.instanceDB.Save(testInstance))

	_, err := sdk.Execute(nil, testInstance.Hash, hash.Int(1), nil, "", testService.Tasks[0].Key, map[string]interface{}{}, nil)
	require.NoError(t, err)

	// not existing instance
	_, err = sdk.Execute(nil, hash.Int(3), hash.Int(1), nil, "", testService.Tasks[0].Key, map[string]interface{}{}, nil)
	require.Error(t, err)
}

func TestExecuteInvalidTaskKey(t *testing.T) {
	sdk, at := newTesting(t)
	defer at.close()

	require.NoError(t, at.serviceDB.Save(testService))

	_, err := sdk.Execute(nil, testService.Hash, hash.Int(1), nil, "", "-", map[string]interface{}{}, nil)
	require.Error(t, err)
}

func TestSubTopic(t *testing.T) {
	require.Equal(t, subTopic(hash.Hash{0}), "1.Execution")
}
