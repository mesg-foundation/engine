package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/app"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xnet"
	"github.com/mesg-foundation/engine/hash"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type apiclient struct {
	pb.EventClient
}

var (
	minExecutionPrice     sdk.Coins
	client                apiclient
	cdc                   = app.MakeCodec()
	processInitialBalance = sdk.NewCoins(sdk.NewInt64Coin("atto", 10000000))
	kb                    *cosmos.Keybase
	cfg                   *config.Config
	engineAddress         sdk.AccAddress
	engineAccountName     string
	engineAccountPassword string
	cont                  container.Container
	ipfsEndpoint          string
	engineName            string
	enginePort            string
	lcd                   *cosmos.LCD
)

const (
	lcdEndpoint     = "http://127.0.0.1:1317/"
	pollingInterval = 500 * time.Millisecond     // half a block
	pollingTimeout  = 10 * time.Second           // 10 blocks
	gasLimit        = flags.DefaultGasLimit * 10 // x10 so the biggest txs have enough gas
)

func TestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// init config
	var err error
	cfg, err = config.New()
	require.NoError(t, err)
	minExecutionPrice, err = sdk.ParseCoins(cfg.DefaultExecutionPrice)
	require.NoError(t, err)
	cosmos.CustomizeConfig(cfg)

	// change and recreate cosmos relative path because CI dir permissions
	cfg.Cosmos.RelativePath = "e2e.cosmos"
	err = os.MkdirAll(filepath.Join(cfg.Path, cfg.Cosmos.RelativePath), os.FileMode(0755))
	require.NoError(t, err)

	// init keybase with account
	kb, err = cosmos.NewKeybase(filepath.Join(cfg.Path, cfg.Cosmos.RelativePath))
	require.NoError(t, err)
	if cfg.Account.Mnemonic != "" {
		acc, err := kb.CreateAccount(cfg.Account.Name, cfg.Account.Mnemonic, "", cfg.Account.Password, keys.CreateHDPath(cfg.Account.Number, cfg.Account.Index).String(), cosmos.DefaultAlgo)
		require.NoError(t, err)
		engineAddress = acc.GetAddress()
		engineAccountName = cfg.Account.Name
		engineAccountPassword = cfg.Account.Password
	}

	// init container
	cont, err = container.New(cfg.Name)
	require.NoError(t, err)
	_, port, _ := xnet.SplitHostPort(cfg.Server.Address)
	enginePort = strconv.Itoa(port)
	engineName = cfg.Name
	ipfsEndpoint = cfg.IpfsEndpoint

	// init gRPC client
	conn, err := grpc.DialContext(context.Background(), "localhost:50052", grpc.WithInsecure())
	require.NoError(t, err)

	client = apiclient{
		pb.NewEventClient(conn),
	}

	// init LCD
	lcd, err = cosmos.NewLCD(lcdEndpoint, cdc, kb, cfg.DevGenesis.ChainID, cfg.Account.Name, cfg.Account.Password, cfg.Cosmos.MinGasPrices, gasLimit)
	require.NoError(t, err)

	// run tests
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

func pollExecution(executionHash hash.Hash, status execution.Status) (*execution.Execution, error) {
	timeout := time.After(pollingTimeout)
	for {
		var exec *execution.Execution
		if err := lcd.Get("execution/get/"+executionHash.String(), &exec); err != nil {
			return nil, err
		}
		if exec.Status == status {
			return exec, nil
		}
		select {
		case <-time.After(pollingInterval):
			continue
		case <-timeout:
			return nil, fmt.Errorf("pollExecution timeout with execution hash %q", executionHash)
		}
	}
}

func pollExecutionOfProcess(processHash hash.Hash, status execution.Status, nodeKey string) (*execution.Execution, error) {
	timeout := time.After(pollingTimeout)
	for {
		var execs []*execution.Execution
		if err := lcd.Get("execution/list", &execs); err != nil {
			return nil, err
		}
		for _, exec := range execs {
			if exec.ProcessHash.Equal(processHash) && exec.Status == status && exec.NodeKey == nodeKey {
				return exec, nil
			}
		}
		select {
		case <-time.After(pollingInterval):
			continue
		case <-timeout:
			return nil, fmt.Errorf("pollExecutionOfProcess timeout with process hash %q and status %q and nodeKey %q", processHash, status, nodeKey)
		}
	}
}
