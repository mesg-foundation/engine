package api

import (
	"context"
	"os"
	"testing"

	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/protobuf/api"
	"github.com/mesg-foundation/core/sdk"
	"github.com/mr-tron/base58"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

const execdbname = "exec.db.test"

func TestGet(t *testing.T) {
	db, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)
	defer db.Close()
	defer os.RemoveAll(execdbname)

	exec := execution.New("", nil, uuid.NewV4().String(), "", nil, nil)
	require.NoError(t, db.Save(exec))

	want, err := toProtoExecution(exec)
	require.NoError(t, err)

	sdk := sdk.New(nil, nil, nil, nil, db)
	s := NewServer(sdk)

	got, err := s.Get(context.Background(), &api.GetExecutionRequest{Hash: base58.Encode(exec.Hash)})
	require.NoError(t, err)
	require.Equal(t, got, want)
}

func TestUpdate(t *testing.T) {
	db, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)
	defer db.Close()
	defer os.RemoveAll(execdbname)

	exec := execution.New("", nil, uuid.NewV4().String(), "", nil, nil)
	require.NoError(t, db.Save(exec))

	sdk := sdk.New(nil, nil, nil, nil, db)
	s := NewServer(sdk)

	_, err = s.Update(context.Background(), &api.UpdateExecutionRequest{Hash: base58.Encode(exec.Hash)})
	require.Equal(t, ErrNoOutput, err)
}
