package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/app"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/cosmos"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type apiclient struct {
	pb.EventClient
	pb.ExecutionClient
	pb.RunnerClient
}

var (
	minExecutionPrice     sdk.Coins
	client                apiclient
	cdc                   = app.MakeCodec()
	processInitialBalance = sdk.NewCoins(sdk.NewInt64Coin("atto", 10000000))
	kb                    *cosmos.Keybase
	cfg                   *config.Config
	engineAddress         sdk.AccAddress
)

func TestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var err error
	cfg, err = config.New()
	require.NoError(t, err)

	minExecutionPrice, err = sdk.ParseCoins(cfg.DefaultExecutionPrice)
	require.NoError(t, err)

	cosmos.CustomizeConfig(cfg)

	conn, err := grpc.DialContext(context.Background(), "localhost:50052", grpc.WithInsecure())
	require.NoError(t, err)

	// change and recreate cosmos relative path because CI dir permissions
	cfg.Cosmos.RelativePath = "e2e.cosmos"
	err = os.MkdirAll(filepath.Join(cfg.Path, cfg.Cosmos.RelativePath), os.FileMode(0755))
	require.NoError(t, err)

	kb, err = cosmos.NewKeybase(filepath.Join(cfg.Path, cfg.Cosmos.RelativePath))
	require.NoError(t, err)
	if cfg.Account.Mnemonic != "" {
		acc, err := kb.CreateAccount(cfg.Account.Name, cfg.Account.Mnemonic, "", cfg.Account.Password, keys.CreateHDPath(cfg.Account.Number, cfg.Account.Index).String(), cosmos.DefaultAlgo)
		require.NoError(t, err)
		engineAddress = acc.GetAddress()
	}

	client = apiclient{
		pb.NewEventClient(conn),
		pb.NewExecutionClient(conn),
		pb.NewRunnerClient(conn),
	}

	// basic tests
	t.Run("service", testService)
	t.Run("runner", testRunner)
	t.Run("process", testProcess)
	t.Run("instance", testInstance)
	t.Run("event", testEvent)
	t.Run("execution", testExecution)
	t.Run("orchestrator", testOrchestrator)
	t.Run("runner/delete", testDeleteRunner)
	t.Run("complex-service", testComplexService)
}
