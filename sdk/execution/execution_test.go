package executionsdk

import (
	"os"
	"testing"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/instance"
	instancesdk "github.com/mesg-foundation/core/sdk/instance"
	servicesdk "github.com/mesg-foundation/core/sdk/service"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

const (
	servicedbname = "service.db.test"
	instdbname    = "instance.db.test"
	execdbname    = "exec.db.test"
)

type apiTesting struct {
	*testing.T
	serviceDB   *database.LevelDBServiceDB
	executionDB *database.LevelDBExecutionDB
	instanceDB  *database.LevelDBInstanceDB
}

func (t *apiTesting) close() {
	require.NoError(t, t.serviceDB.Close())
	require.NoError(t, t.executionDB.Close())
	require.NoError(t, t.instanceDB.Close())
	require.NoError(t, os.RemoveAll(servicedbname))
	require.NoError(t, os.RemoveAll(execdbname))
	require.NoError(t, os.RemoveAll(instdbname))
}

func newTesting(t *testing.T) (*Execution, *apiTesting) {
	container := &mocks.Container{}
	db, err := database.NewServiceDB(servicedbname)
	require.NoError(t, err)
	service := servicesdk.New(container, db)

	instDB, err := database.NewInstanceDB(instdbname)
	require.NoError(t, err)
	instance := instancesdk.New(container, service, instDB)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	sdk := New(pubsub.New(0), service, instance, execDB)

	return sdk, &apiTesting{
		T:           t,
		serviceDB:   db,
		executionDB: execDB,
		instanceDB:  instDB,
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
	exec := execution.New(nil, nil, "", "", nil, nil)
	require.NoError(t, sdk.execDB.Save(exec))
	got, err := sdk.Get(exec.Hash)
	require.NoError(t, err)
	require.Equal(t, exec, got)
}

func TestGetStream(t *testing.T) {
	sdk, at := newTesting(t)
	defer at.close()

	exec := execution.New(nil, nil, "", "", nil, nil)
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

	_, err := sdk.Execute(testInstance.Hash, testService.Tasks[0].Key, map[string]interface{}{}, []string{})
	require.NoError(t, err)
}

func TestExecuteWithInvalidTaskName(t *testing.T) {
	sdk, at := newTesting(t)
	defer at.close()

	require.NoError(t, at.serviceDB.Save(testService))

	_, err := sdk.Execute(testService.Hash, "-", nil, nil)
	require.Error(t, err)
}

func TestSubTopic(t *testing.T) {
	require.Equal(t, subTopic(hash.Hash{0}), "1.Execution")
}
