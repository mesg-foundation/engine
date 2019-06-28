package api

import (
	"context"
	"os"
	"testing"

	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/sdk"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

const execdbname = "exec.db.test"

func TestGet(t *testing.T) {
	db, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)
	defer db.Close()
	defer os.RemoveAll(execdbname)

	exec := execution.New(nil, nil, uuid.NewV4().String(), "", nil, nil)
	require.NoError(t, db.Save(exec))

	want, err := toProtoExecution(exec)
	require.NoError(t, err)

	sdk := sdk.New(nil, nil, nil, db)
	s := NewExecutionServer(sdk)

	got, err := s.Get(context.Background(), &api.GetExecutionRequest{Hash: exec.Hash.String()})
	require.NoError(t, err)
	require.Equal(t, got, want)
}

func TestUpdate(t *testing.T) {
	db, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)
	defer db.Close()
	defer os.RemoveAll(execdbname)

	exec := execution.New(nil, nil, uuid.NewV4().String(), "", nil, nil)
	require.NoError(t, db.Save(exec))

	sdk := sdk.New(nil, nil, nil, db)
	s := NewExecutionServer(sdk)

	_, err = s.Update(context.Background(), &api.UpdateExecutionRequest{Hash: exec.Hash.String()})
	require.Equal(t, ErrNoOutput, err)
}
