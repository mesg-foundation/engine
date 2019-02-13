package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/ipfs"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/utils/readme"
	uuid "github.com/satori/go.uuid"
)

// MarketplaceProvider is a struct that provides all methods required by service command.
type MarketplaceProvider struct {
	client coreapi.CoreClient
}

// NewMarketplaceProvider creates new MarketplaceProvider.
func NewMarketplaceProvider(c coreapi.CoreClient) *MarketplaceProvider {
	return &MarketplaceProvider{client: c}
}

const (
	// PublishVersion is the version used to publish the services to the marketplace
	PublishVersion = 1

	// DeploymentType is the type of deployment used for the service
	DeploymentType = "ipfs"

	// MarketplaceServiceID is the sid of the marketplace service
	MarketplaceServiceID = "marketplace"
)

type ManifestData struct {
	Version int `json:"version"`
	Service struct {
		Deployment struct {
			Type   string `json:"type"`
			Source string `json:"source"`
		} `json:"deployment"`
	} `json:"service"`
	Definition *importer.ServiceDefinition `json:"definition"`
	Readme     string                      `json:"readme,omitempty"`
}

// TransactionTaskInputs is the inputs for any task that create a transaction.
type TransactionTaskInputs struct {
	From string `json:"from"`
	// Gas      string `json:"gas"` // omitempty
	// GasPrice string `json:"gasPrice"` // omitempty
}

// ErrorOutput is the output for any task that fails.
type ErrorOutput struct {
	Message string `json:"message"`
}

// TransactionOutput is the output for any task that creates a transaction.
type TransactionOutput struct {
	ChainID  int64  `json:"chainID"`
	Nonce    uint64 `json:"nonce"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Gas      uint64 `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Data     string `json:"data"`
}

// CreateServiceTaskInputs is the inputs of the task create service.
type CreateServiceTaskInputs struct {
	*TransactionTaskInputs
	Sid string `json:"sid"`
}

// CreateServiceVersionTaskInputs is the inputs of the task create service version.
type CreateServiceVersionTaskInputs struct {
	*TransactionTaskInputs
	SidHash          string `json:"sidHash"`
	Hash             string `json:"hash"`
	Manifest         string `json:"manifest"`
	ManifestProtocol string `json:"manifestProtocol"`
}

// CreateServiceOfferTaskInputs is the inputs of the task create service offer.
type CreateServiceOfferTaskInputs struct {
	*TransactionTaskInputs
	SidHash  string `json:"sidHash"`
	Price    string `json:"price"`
	Duration string `json:"duration"`
}

// DisableServiceOfferTaskInputs is the inputs of the task create service offer.
type DisableServiceOfferTaskInputs struct {
	*TransactionTaskInputs
	SidHash    string `json:"sidHash"`
	OfferIndex string `json:"offerIndex"`
}

// PurchaseTaskInputs is the inputs of the task purchase.
type PurchaseTaskInputs struct {
	*TransactionTaskInputs
	SidHash    string `json:"sidHash"`
	OfferIndex string `json:"offerIndex"`
}

// TransferServiceOwnershipTaskInputs is the inputs of the task transfer service ownership.
type TransferServiceOwnershipTaskInputs struct {
	*TransactionTaskInputs
	SidHash  string `json:"sidHash"`
	NewOwner string `json:"newOwner"`
}

// SendSignedTransactionTaskInputs is the inputs of the task send signed transaction.
type SendSignedTransactionTaskInputs struct {
	SignedTransaction string `json:"signedTransaction"`
}

// SendSignedTransactionTaskSuccessOutput is the success output of task send signed transaction.
type SendSignedTransactionTaskSuccessOutput struct {
	Receipt string `json:"receipt"`
}

// IsAuthorizedTaskInputs is the inputs of the task is authorized.
type IsAuthorizedTaskInputs struct {
	SidHash string `json:"sidHash"`
}

// IsAuthorizedTaskSuccessOutput is the success output of task authorized.
type IsAuthorizedTaskSuccessOutput struct {
	Authorized bool `json:"authorized"`
}

// CreateService executes the create service task
func (p *MarketplaceProvider) CreateService(sid, from string) (*TransactionOutput, error) {
	outputKey, outputDataS, err := p.executeTaskAndWaitResult(MarketplaceServiceID, "createService", &CreateServiceTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{
			From: from,
		},
		Sid: sid,
	})
	if err != nil {
		return nil, err
	}
	var outputData *TransactionOutput
	return outputData, decodeOutput(outputKey, outputDataS, outputData)
}

// CreateServiceVersion executes the create service version task
func (p *MarketplaceProvider) CreateServiceVersion(sidHash, hash, manifest, manifestProtocol, from string) (*TransactionOutput, error) {
	outputKey, outputDataS, err := p.executeTaskAndWaitResult(MarketplaceServiceID, "createServiceVersion", &CreateServiceVersionTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{
			From: from,
		},
		SidHash:          sidHash,
		Hash:             hash,
		Manifest:         manifest,
		ManifestProtocol: manifestProtocol,
	})
	if err != nil {
		return nil, err
	}
	var outputData *TransactionOutput
	return outputData, decodeOutput(outputKey, outputDataS, outputData)
}

// CreateServiceOffer executes the create service offer task
func (p *MarketplaceProvider) CreateServiceOffer(sidHash, price, duration, from string) (*TransactionOutput, error) {
	outputKey, outputDataS, err := p.executeTaskAndWaitResult(MarketplaceServiceID, "createServiceOffer", &CreateServiceOfferTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{
			From: from,
		},
		SidHash:  sidHash,
		Price:    price,
		Duration: duration,
	})
	if err != nil {
		return nil, err
	}
	var outputData *TransactionOutput
	return outputData, decodeOutput(outputKey, outputDataS, outputData)
}

// DisableServiceOffer executes the disable service offer task
func (p *MarketplaceProvider) DisableServiceOffer(sidHash, offerIndex, from string) (*TransactionOutput, error) {
	outputKey, outputDataS, err := p.executeTaskAndWaitResult(MarketplaceServiceID, "disableServiceOffer", &DisableServiceOfferTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{
			From: from,
		},
		SidHash:    sidHash,
		OfferIndex: offerIndex,
	})
	if err != nil {
		return nil, err
	}
	var outputData *TransactionOutput
	return outputData, decodeOutput(outputKey, outputDataS, outputData)
}

// Purchase executes the purchase task
func (p *MarketplaceProvider) Purchase(sidHash, offerIndex, from string) (*TransactionOutput, error) {
	outputKey, outputDataS, err := p.executeTaskAndWaitResult(MarketplaceServiceID, "purchase", &PurchaseTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{
			From: from,
		},
		SidHash:    sidHash,
		OfferIndex: offerIndex,
	})
	if err != nil {
		return nil, err
	}
	var outputData *TransactionOutput
	return outputData, decodeOutput(outputKey, outputDataS, outputData)
}

// TransferServiceOwnership executes the task transfer service ownership.
func (p *MarketplaceProvider) TransferServiceOwnership(sidHash, newOwner, from string) (*TransactionOutput, error) {
	outputKey, outputDataS, err := p.executeTaskAndWaitResult(MarketplaceServiceID, "transferServiceOwnership", &TransferServiceOwnershipTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{
			From: from,
		},
		SidHash:  sidHash,
		NewOwner: newOwner,
	})
	if err != nil {
		return nil, err
	}
	var outputData *TransactionOutput
	return outputData, decodeOutput(outputKey, outputDataS, outputData)
}

// SendSignedTransaction executes the task send signed transaction.
func (p *MarketplaceProvider) SendSignedTransaction(signedTransaction string) (*SendSignedTransactionTaskSuccessOutput, error) {
	outputKey, outputDataS, err := p.executeTaskAndWaitResult(MarketplaceServiceID, "sendSignedTransaction", &SendSignedTransactionTaskInputs{
		SignedTransaction: signedTransaction,
	})
	if err != nil {
		return nil, err
	}
	var outputData *SendSignedTransactionTaskSuccessOutput
	return outputData, decodeOutput(outputKey, outputDataS, outputData)
}

// IsAuthorized executes the task send signed transaction.
func (p *MarketplaceProvider) IsAuthorized(sidHash string) (bool, error) {
	outputKey, outputDataS, err := p.executeTaskAndWaitResult(MarketplaceServiceID, "isAuthorized", &IsAuthorizedTaskInputs{
		SidHash: sidHash,
	})
	if err != nil {
		return false, err
	}
	var outputData *IsAuthorizedTaskSuccessOutput
	if err = decodeOutput(outputKey, outputDataS, outputData); err != nil {
		return false, err
	}
	return outputData.Authorized, nil
}

// PublishDefinitionFile upload and publish the tarball and definition file and returns the address of the definition file
func (p *MarketplaceProvider) PublishDefinitionFile(path string) (string, error) {
	ipfs := ipfs.New()
	tar, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return "", err
	}
	tarballResponse, err := ipfs.Add("tarball", tar)
	if err != nil {
		return "", err
	}

	definitionFile, err := p.createDefinitionFile(path, tarballResponse.Hash)
	if err != nil {
		return "", err
	}

	definitionResponse, err := ipfs.Add("definition", bytes.NewReader(definitionFile))
	if err != nil {
		return "", err
	}
	return definitionResponse.Hash, nil
}

func (p *MarketplaceProvider) createDefinitionFile(path string, tarballHash string) ([]byte, error) {
	definition, err := importer.From(path)
	if err != nil {
		return nil, err
	}
	var data ManifestData
	data.Version = PublishVersion
	data.Service.Deployment.Type = DeploymentType
	data.Service.Deployment.Source = tarballHash
	data.Readme, err = readme.Lookup(path)
	if err != nil {
		return nil, err
	}
	data.Definition = definition
	return json.Marshal(data)
}

// TODO: to refactor with the same function in service_provider
func (p *MarketplaceProvider) executeTaskAndWaitResult(sid, taskKey string, inputData interface{}) (string, string, error) {
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
func (p *MarketplaceProvider) listenEvents(id, eventFilter string) (chan *coreapi.EventData, chan error, error) {
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
func (p *MarketplaceProvider) listenResults(id, taskFilter, outputFilter string, tagFilters []string) (chan *coreapi.ResultData, chan error, error) {
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
func (p *MarketplaceProvider) executeTask(id, taskKey, inputData string, tags []string) error {
	_, err := p.client.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID:     id,
		TaskKey:       taskKey,
		InputData:     inputData,
		ExecutionTags: tags,
	})
	return err
}

func decodeOutput(outputKey, outputDataS string, outputData interface{}) error {
	if outputKey != "success" {
		return fmt.Errorf("Task returns an error: %s", outputDataS)
	}
	return json.Unmarshal([]byte(outputDataS), &outputData)
}
