package api

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func newAPI(t *testing.T) (*API, *service.Service, func()) {
	name, err := ioutil.TempDir("", "db")
	require.NoError(t, err)
	db, err := database.NewServiceDB(filepath.Join(name, "service"))
	require.NoError(t, err)
	execdb, err := database.NewExecutionDB(filepath.Join(name, "exec"))
	require.NoError(t, err)

	s := &service.Service{
		Hash:   "0",
		Events: []*service.Event{{Key: "foo"}},
		Tasks:  []*service.Task{{Key: "foo"}},
	}
	require.NoError(t, db.Save(s))

	return New(db, execdb, service.NewFakeManager()), s, func() {
		db.Close()
		execdb.Close()
		os.RemoveAll(name)
	}
}

func TestGetService(t *testing.T) {
	api, s, closer := newAPI(t)
	defer closer()

	ts, err := api.GetService(s.Hash)
	require.NoError(t, err)
	require.Equal(t, s.Hash, ts.Hash)
	require.Equal(t, service.StatusUnknown, ts.Status)
}

func TestDeleteService(t *testing.T) {
	api, s, closer := newAPI(t)
	defer closer()

	require.NoError(t, api.DeleteService(s.Hash, true))
	api.sm.Status(s)
	require.Equal(t, service.StatusDeleted, s.Status)
}

func TestServiceLogs(t *testing.T) {
	api, s, closer := newAPI(t)
	defer closer()

	logs, err := api.ServiceLogs(s.Hash, nil)
	require.NoError(t, err)
	require.Len(t, logs, 1)
}

func TestListService(t *testing.T) {
	api, s, closer := newAPI(t)
	defer closer()

	services, err := api.ListServices()
	require.NoError(t, err)
	require.Len(t, services, 1)
	require.Equal(t, s.Hash, services[0].Hash)
}

func TestStopService(t *testing.T) {
	api, s, closer := newAPI(t)
	defer closer()

	require.NoError(t, api.StopService(s.Hash))
	api.sm.Status(s)
	require.Equal(t, service.StatusStopped, s.Status)
}

func TestStartService(t *testing.T) {
	api, s, closer := newAPI(t)
	defer closer()

	require.NoError(t, api.StartService(s.Hash))
	api.sm.Status(s)
	require.Equal(t, service.StatusRunning, s.Status)
}

func TestEmitEvent(t *testing.T) {
	api, s, closer := newAPI(t)
	defer closer()

	channel := s.EventSubscriptionChannel()
	ch := pubsub.Subscribe(channel)
	defer pubsub.Unsubscribe(channel, ch)

	require.NoError(t, api.EmitEvent(s.Hash, "foo", nil))
	event := (<-ch).(*event.Event)
	require.Equal(t, s.Hash, event.Service.Hash)
}

func TestExecuteTask(t *testing.T) {
	api, s, closer := newAPI(t)
	defer closer()

	require.NoError(t, api.StartService(s.Hash))

	channel := s.TaskSubscriptionChannel()
	ch := pubsub.Subscribe(channel)
	defer pubsub.Unsubscribe(channel, ch)

	id, err := api.ExecuteTask(s.Hash, "foo", nil, nil)
	require.NoError(t, err)
	exec := (<-ch).(*execution.Execution)
	require.Equal(t, s.Hash, exec.Service.Hash)
	require.Equal(t, id, exec.ID)
}
