package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/app"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type apiclient struct {
	EventClient     orchestrator.EventClient
	ExecutionClient orchestrator.ExecutionClient
	RunnerClient    orchestrator.RunnerClient
}

const (
	chainID         = "mesg-dev-chain"
	gasPrices       = ""
	pollingInterval = 500 * time.Millisecond // half a block
	pollingTimeout  = 10 * time.Second       // 10 blocks

	lcdEndpoint          = "http://127.0.0.1:1317/"
	orchestratorEndpoint = "localhost:50052"
	orchestratorName     = "engine"
	orchestratorNetwork  = "engine"

	engineMnemonic        = "neutral false together tattoo matrix stamp poem mouse chair chair grain pledge mandate layer shiver embark struggle vicious antenna total faith genre valley mandate"
	engineAccountName     = "engine"
	engineAccountPassword = "pass"
	engineAccountNumber   = uint32(0)
	engineAccountIndex    = uint32(0)

	cliAccountMnemonic = "spike raccoon obscure program raw large unaware dragon hamster round artist case fall wage sample velvet robust legend identify innocent film coral picture organ"
	cliAccountName     = "cli"
	cliAccountPassword = "pass"
	cliAccountNumber   = uint32(0)
	cliAccountIndex    = uint32(0)
)

var (
	executionPrice        = sdk.NewCoins(sdk.NewInt64Coin("atto", 10000)) // /x/execution/internal/types/params.go#DefaultMinPrice
	processInitialBalance = sdk.NewCoins(sdk.NewInt64Coin("atto", 10000000))

	cdc           = app.MakeCodec()
	client        *apiclient
	kb            *cosmos.Keybase
	engineAddress sdk.AccAddress
	cont          *container.Container
	lcd           *cosmos.LCD
	lcdEngine     *cosmos.LCD
	cliAddress    sdk.AccAddress
)

func TestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// init the config of cosmos
	cosmos.InitConfig()

	// init config
	var err error

	// init keybase with engine account and cli account
	kb = cosmos.NewInMemoryKeybase()

	// init engine account
	engineAcc, err := kb.CreateAccount(engineAccountName, engineMnemonic, "", engineAccountPassword, keys.CreateHDPath(engineAccountNumber, engineAccountIndex).String(), cosmos.DefaultAlgo)
	require.NoError(t, err)
	engineAddress = engineAcc.GetAddress()

	// init cli account
	cliAcc, err := kb.CreateAccount(cliAccountName, cliAccountMnemonic, "", cliAccountPassword, keys.CreateHDPath(cliAccountNumber, cliAccountIndex).String(), cosmos.DefaultAlgo)
	require.NoError(t, err)
	cliAddress = cliAcc.GetAddress()

	// init LCD with engine account and make a transfer to cli account
	lcdEngine, err = cosmos.NewLCD(lcdEndpoint, cdc, kb, chainID, engineAccountName, engineAccountPassword, gasPrices)
	require.NoError(t, err)

	// init container
	cont, err = container.New(orchestratorName, orchestratorEndpoint, orchestratorNetwork, 0, 5*time.Second)
	require.NoError(t, err)

	// init orchestrator gRPC client
	conn, err := grpc.DialContext(context.Background(), orchestratorEndpoint, grpc.WithInsecure())
	require.NoError(t, err)

	client = &apiclient{
		EventClient:     orchestrator.NewEventClient(conn),
		ExecutionClient: orchestrator.NewExecutionClient(conn),
		RunnerClient:    orchestrator.NewRunnerClient(conn),
	}

	// init LCD
	lcd, err = cosmos.NewLCD(lcdEndpoint, cdc, kb, chainID, cliAccountName, cliAccountPassword, gasPrices)
	require.NoError(t, err)

	// run tests
	t.Run("account-sequence", testAccountSequence)
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
		if err := lcd.Get(fmt.Sprintf("execution/list?processHash=%s&status=%s&nodeKey=%s", processHash, status, nodeKey), &execs); err != nil {
			return nil, err
		}
		if len(execs) > 0 {
			return execs[0], nil
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
