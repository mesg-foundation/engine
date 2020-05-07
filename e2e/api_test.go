package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/mesg-foundation/engine/app"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type apiclient struct {
	EventClient        orchestrator.EventClient
	ExecutionClient    orchestrator.ExecutionClient
	RunnerClient       orchestrator.RunnerClient
	OrchestratorClient orchestrator.OrchestratorClient
}

var (
	minExecutionPrice     sdk.Coins
	client                *apiclient
	cdc                   = app.MakeCodec()
	processInitialBalance = sdk.NewCoins(sdk.NewInt64Coin("atto", 10000000))
	kb                    *cosmos.Keybase
	cfg                   *config.Config
	engineAddress         sdk.AccAddress
	cont                  *container.Container
	lcd                   *cosmos.LCD
	lcdEngine             *cosmos.LCD
	cliAddress            sdk.AccAddress
	cliInitialBalance, _  = sdk.ParseCoins("100000000000000000000000000atto")
)

const (
	lcdEndpoint        = "http://127.0.0.1:1317/"
	pollingInterval    = 500 * time.Millisecond // half a block
	pollingTimeout     = 10 * time.Second       // 10 blocks
	cliAccountMnemonic = "large fork soccer lab answer enlist robust vacant narrow please inmate primary father must add hub shy couch rail video tool marine pill give"
	cliAccountName     = "cli"
	cliAccountPassword = "pass"
)

func TestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// init the config of cosmos
	cosmos.InitConfig()

	// init config
	var err error
	cfg, err = config.New()
	require.NoError(t, err)
	minExecutionPrice, err = sdk.ParseCoins(cfg.DefaultExecutionPrice)
	require.NoError(t, err)

	// change and recreate cosmos relative path because CI dir permissions
	cfg.Cosmos.RelativePath = "e2e.cosmos"
	err = os.MkdirAll(filepath.Join(cfg.Path, cfg.Cosmos.RelativePath), os.FileMode(0755))
	require.NoError(t, err)

	// init keybase with engine account and cli account
	kb, err = cosmos.NewKeybase(filepath.Join(cfg.Path, cfg.Cosmos.RelativePath))
	require.NoError(t, err)
	// init engine account
	engineAcc, err := kb.CreateAccount(cfg.Account.Name, cfg.Account.Mnemonic, "", cfg.Account.Password, keys.CreateHDPath(cfg.Account.Number, cfg.Account.Index).String(), cosmos.DefaultAlgo)
	require.NoError(t, err)
	engineAddress = engineAcc.GetAddress()

	// init cli account
	cliAcc, err := kb.CreateAccount(cliAccountName, cliAccountMnemonic, "", cliAccountPassword, keys.CreateHDPath(cfg.Account.Number, cfg.Account.Index).String(), cosmos.DefaultAlgo)
	require.NoError(t, err)
	cliAddress = cliAcc.GetAddress()

	// init LCD with engine account and make a transfer to cli account
	lcdEngine, err = cosmos.NewLCD(lcdEndpoint, cdc, kb, cfg.DevGenesis.ChainID, cfg.Account.Name, cfg.Account.Password, cfg.Cosmos.MinGasPrices)
	require.NoError(t, err)
	_, err = lcdEngine.BroadcastMsg(bank.NewMsgSend(engineAddress, cliAddress, cliInitialBalance))
	require.NoError(t, err)

	// init container
	cont, err = container.New(cfg.Name, cfg.Server.Address, cfg.Name)
	require.NoError(t, err)

	// init gRPC client
	conn, err := grpc.DialContext(context.Background(), "localhost:50052", grpc.WithInsecure())
	require.NoError(t, err)

	client = &apiclient{
		EventClient:        orchestrator.NewEventClient(conn),
		ExecutionClient:    orchestrator.NewExecutionClient(conn),
		RunnerClient:       orchestrator.NewRunnerClient(conn),
		OrchestratorClient: orchestrator.NewOrchestratorClient(conn),
	}

	// init LCD
	lcd, err = cosmos.NewLCD(lcdEndpoint, cdc, kb, cfg.DevGenesis.ChainID, cliAccountName, cliAccountPassword, cfg.Cosmos.MinGasPrices)
	require.NoError(t, err)

	// run tests
	// t.Run("account-sequence", testAccountSequence)
	t.Run("service", testService)
	t.Run("runner", testRunner)
	// t.Run("process", testProcess)
	// t.Run("instance", testInstance)
	// t.Run("event", testEvent)
	// t.Run("execution", testExecution)
	t.Run("orchestrator", testOrchestrator)
	t.Run("runner/delete", testDeleteRunner)
	// t.Run("complex-service", testComplexService)
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

func signPayload(payload interface{}) ([]byte, error) {
	encodedValue, err := cdc.MarshalJSON(payload)
	if err != nil {
		return nil, err
	}
	encodedValueSorted, err := sdk.SortJSON(encodedValue)
	if err != nil {
		return nil, err
	}
	signature, _, err := kb.Sign(cliAccountName, cliAccountPassword, encodedValueSorted)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// signCred is a structure that manage a token.
type signCred struct {
	request interface{}
}

// GetRequestMetadata returns the metadata for the request.
func (c *signCred) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	signature, err := signPayload(c.request)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		orchestrator.RequestSignature: base64.StdEncoding.EncodeToString(signature),
	}, nil
}

// RequireTransportSecurity tells if the transport should be secured.
func (c *signCred) RequireTransportSecurity() bool {
	return false
}
