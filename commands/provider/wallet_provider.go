package provider

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	uuid "github.com/satori/go.uuid"
)

// WalletProvider is a struct that provides all methods required by wallet command.
type WalletProvider struct {
	client coreapi.CoreClient
}

// NewWalletProvider creates new WalletProvider.
func NewWalletProvider(client coreapi.CoreClient) *WalletProvider {
	return &WalletProvider{client: client}
}

func (p *WalletProvider) List() ([]common.Address, error) { return nil, nil }
func (p *WalletProvider) Create(passphrase string) (common.Address, error) {
	return common.Address{}, nil
}
func (p *WalletProvider) Delete(address common.Address, passphrase string) error { return nil }
func (p *WalletProvider) Export(address common.Address, passphrase string) ([]byte, error) {
	return nil, nil
}
func (p *WalletProvider) Import(address common.Address, passphrase string, account []byte) error {
	return nil
}

const WalletServiceID = "ethwallet"

// SignTaskInputs is the inputs of the task sign.
type SignTaskInputs struct {
	*TransactionTaskInputs
	Address     string             `json:"address"`
	Passphrase  string             `json:"passphrase"`
	Transaction *TransactionOutput `json:"transaction"`
}

// SignTaskSuccessOutput is the success output of task authorized.
type SignTaskSuccessOutput struct {
	SignedTransaction string `json:"signedTransaction"`
}

func (p *WalletProvider) Sign(address string, passphrase string, transaction *TransactionOutput) (string, error) {
	outputKey, outputDataS, err := p.executeTaskAndWaitResult(WalletServiceID, "sign", &SignTaskInputs{
		Address:     address,
		Passphrase:  passphrase,
		Transaction: transaction,
	})
	if err != nil {
		return "", err
	}
	var outputData SignTaskSuccessOutput
	if err = decodeOutput(outputKey, outputDataS, &outputData); err != nil {
		return "", err
	}
	return outputData.SignedTransaction, nil
}

// TODO: to refactor with the same function in service_provider
func (p *WalletProvider) executeTaskAndWaitResult(sid, taskKey string, inputData interface{}) (string, string, error) {
	var (
		result *coreapi.ResultData
		err    error
		data   []byte
		tags   = []string{uuid.NewV4().String()}
	)
	data, err = json.Marshal(inputData)
	if err != nil {
		return "", "", err
	}
	// TODO: the stream should be terminated when one result received
	listenResultsC, resultsErrC, err := p.listenResults(sid, taskKey, "", tags)
	if err != nil {
		return "", "", err
	}
	err = p.executeTask(sid, taskKey, string(data), tags)
	if err != nil {
		return "", "", err
	}
	select {
	case result = <-listenResultsC:
	case err = <-resultsErrC:
	}
	if err != nil {
		return "", "", err
	}
	if result.Error != "" {
		return "", "", errors.New(result.Error)
	}
	return result.OutputKey, result.OutputData, nil
}

// listenEvents returns a channel with event data streaming.
func (p *WalletProvider) listenEvents(id, eventFilter string) (chan *coreapi.EventData, chan error, error) {
	stream, err := p.client.ListenEvent(context.Background(), &coreapi.ListenEventRequest{
		ServiceID:   id,
		EventFilter: eventFilter,
	})
	if err != nil {
		return nil, nil, err
	}

	resultC := make(chan *coreapi.EventData)
	errC := make(chan error)

	go func() {
		<-stream.Context().Done()
		errC <- stream.Context().Err()
	}()
	go func() {
		for {
			if res, err := stream.Recv(); err != nil {
				errC <- err
				break
			} else {
				resultC <- res
			}
		}
	}()

	if err := acknowledgement.WaitForStreamToBeReady(stream); err != nil {
		return nil, nil, err
	}

	return resultC, errC, nil
}

// listenResults returns a channel with event results streaming..
func (p *WalletProvider) listenResults(id, taskFilter, outputFilter string, tagFilters []string) (chan *coreapi.ResultData, chan error, error) {
	stream, err := p.client.ListenResult(context.Background(), &coreapi.ListenResultRequest{
		ServiceID:    id,
		TaskFilter:   taskFilter,
		OutputFilter: outputFilter,
		TagFilters:   tagFilters,
	})
	if err != nil {
		return nil, nil, err
	}
	resultC := make(chan *coreapi.ResultData)
	errC := make(chan error)

	go func() {
		<-stream.Context().Done()
		errC <- stream.Context().Err()
	}()
	go func() {
		for {
			if res, err := stream.Recv(); err != nil {
				errC <- err
				break
			} else {
				resultC <- res
			}
		}
	}()

	if err := acknowledgement.WaitForStreamToBeReady(stream); err != nil {
		return nil, nil, err
	}

	return resultC, errC, nil
}

// executeTask executes task on given service.
func (p *WalletProvider) executeTask(id, taskKey, inputData string, tags []string) error {
	_, err := p.client.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID:     id,
		TaskKey:       taskKey,
		InputData:     inputData,
		ExecutionTags: tags,
	})
	return err
}
