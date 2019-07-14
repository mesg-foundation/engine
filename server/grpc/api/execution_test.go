package api

import (
	"context"
	"os"
	"testing"

	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
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

	exec := execution.New(nil, nil, nil, "", nil, nil)
	require.NoError(t, db.Save(exec))

	want, err := toProtoExecution(exec)
	require.NoError(t, err)

	sdk := sdk.New(nil, nil, nil, db, "", "")
	s := NewExecutionServer(sdk)

	got, err := s.Get(context.Background(), &api.GetExecutionRequest{Hash: exec.Hash.String()})
	require.NoError(t, err)
	require.Equal(t, got, want)
}

func TestCreateWithInvalidHash(t *testing.T) {
	sdk := sdk.New(nil, nil, nil, nil, "", "")
	s := NewExecutionServer(sdk)

	_, err := s.Create(context.Background(), &api.CreateExecutionRequest{
		InstanceHash: "xx",
	})
	require.Error(t, err)
}

func TestCreateWithInvalidInputs(t *testing.T) {
	sdk := sdk.New(nil, nil, nil, nil, "", "")
	s := NewExecutionServer(sdk)

	h, _ := hash.Random()
	_, err := s.Create(context.Background(), &api.CreateExecutionRequest{
		InstanceHash: h.String(),
		Inputs:       "xx",
	})
	require.Error(t, err)
}

func TestUpdate(t *testing.T) {
	db, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)
	defer db.Close()
	defer os.RemoveAll(execdbname)

	exec := execution.New(nil, nil, nil, "", nil, nil)
	require.NoError(t, db.Save(exec))

	sdk := sdk.New(nil, nil, nil, db, "", "")
	s := NewExecutionServer(sdk)

	_, err = s.Update(context.Background(), &api.UpdateExecutionRequest{Hash: exec.Hash.String()})
	require.Equal(t, ErrNoOutput, err)
}
