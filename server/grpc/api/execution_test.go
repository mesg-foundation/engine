package api

import (
	"context"
	"os"
	"testing"

	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/sdk"
	"github.com/stretchr/testify/require"
)

const execdbname = "exec.db.test"

func TestGet(t *testing.T) {
	db, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)
	defer db.Close()
	defer os.RemoveAll(execdbname)

	exec := execution.New(nil, nil, nil, nil, "", "", nil, nil)
	require.NoError(t, db.Save(exec))

	sdk := sdk.New(nil, nil, nil, db, nil, nil, "", "", "")
	s := NewExecutionServer(sdk)

	got, err := s.Get(context.Background(), &api.GetExecutionRequest{Hash: exec.Hash})
	require.NoError(t, err)
	require.True(t, got.Equal(exec))
}

func TestUpdate(t *testing.T) {
	db, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)
	defer db.Close()
	defer os.RemoveAll(execdbname)

	exec := execution.New(nil, nil, nil, nil, "", "", nil, nil)
	require.NoError(t, db.Save(exec))

	sdk := sdk.New(nil, nil, nil, db, nil, nil, "", "", "")
	s := NewExecutionServer(sdk)

	_, err = s.Update(context.Background(), &api.UpdateExecutionRequest{Hash: exec.Hash})
	require.Equal(t, ErrNoOutput, err)
}
