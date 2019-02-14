package provider

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/ipfs"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/utils/readme"
)

// MarketplaceProvider is a struct that provides all methods required by service command.
type MarketplaceProvider struct {
	client client
}

// NewMarketplaceProvider creates new MarketplaceProvider.
func NewMarketplaceProvider(c coreapi.CoreClient) *MarketplaceProvider {
	return &MarketplaceProvider{client: client{c}}
}

// CreateService executes the create service task
func (p *MarketplaceProvider) CreateService(sid, from string) (*Transaction, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "createService", &CreateServiceTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{
			From: from,
		},
		Sid: sid,
	})
	if err != nil {
		return nil, err
	}

	if r.OutputKey == "error" {
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return nil, err
		}
		return nil, errors.New(output.Message)
	}

	var output Transaction
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// CreateServiceVersion executes the create service version task
func (p *MarketplaceProvider) CreateServiceVersion(sidHash, hash, manifest, manifestProtocol, from string) (*Transaction, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "createServiceVersion", &CreateServiceVersionTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{From: from},
		SidHash:               sidHash,
		Hash:                  hash,
		Manifest:              manifest,
		ManifestProtocol:      manifestProtocol,
	})
	if err != nil {
		return nil, err
	}

	if r.OutputKey == "error" {
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return nil, err
		}
		return nil, errors.New(output.Message)
	}

	var output Transaction
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// CreateServiceOffer executes the create service offer task
func (p *MarketplaceProvider) CreateServiceOffer(sidHash, price, duration, from string) (*Transaction, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "createServiceOffer", &CreateServiceOfferTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{From: from},
		SidHash:               sidHash,
		Price:                 price,
		Duration:              duration,
	})
	if err != nil {
		return nil, err
	}

	if r.OutputKey == "error" {
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return nil, err
		}
		return nil, errors.New(output.Message)
	}

	var output Transaction
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// DisableServiceOffer executes the disable service offer task
func (p *MarketplaceProvider) DisableServiceOffer(sidHash, offerIndex, from string) (*Transaction, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "disableServiceOffer", &DisableServiceOfferTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{From: from},
		SidHash:               sidHash,
		OfferIndex:            offerIndex,
	})
	if err != nil {
		return nil, err
	}

	if r.OutputKey == "error" {
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return nil, err
		}
		return nil, errors.New(output.Message)
	}

	var output Transaction
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// Purchase executes the purchase task
func (p *MarketplaceProvider) Purchase(sidHash, offerIndex, from string) (*Transaction, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "purchase", &PurchaseTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{From: from},
		SidHash:               sidHash,
		OfferIndex:            offerIndex,
	})
	if err != nil {
		return nil, err
	}

	if r.OutputKey == "error" {
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return nil, err
		}
		return nil, errors.New(output.Message)
	}

	var output Transaction
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// TransferServiceOwnership executes the task transfer service ownership.
func (p *MarketplaceProvider) TransferServiceOwnership(sidHash, newOwner, from string) (*Transaction, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "transferServiceOwnership", &TransferServiceOwnershipTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{From: from},
		SidHash:               sidHash,
		NewOwner:              newOwner,
	})
	if err != nil {
		return nil, err
	}

	if r.OutputKey == "error" {
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return nil, err
		}
		return nil, errors.New(output.Message)
	}

	var output Transaction
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// SendSignedTransaction executes the task send signed transaction.
func (p *MarketplaceProvider) SendSignedTransaction(signedTransaction string) (*SendSignedTransactionTaskSuccessOutput, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "sendSignedTransaction", &SendSignedTransactionTaskInputs{
		SignedTransaction: signedTransaction,
	})
	if err != nil {
		return nil, err
	}

	if r.OutputKey == "error" {
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return nil, err
		}
		return nil, errors.New(output.Message)
	}

	var output SendSignedTransactionTaskSuccessOutput
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// IsAuthorized executes the task send signed transaction.
func (p *MarketplaceProvider) IsAuthorized(sidHash string) (bool, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "isAuthorized", &IsAuthorizedTaskInputs{
		SidHash: sidHash,
	})
	if err != nil {
		return false, err
	}

	if r.OutputKey == "error" {
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return false, err
		}
		return false, errors.New(output.Message)
	}

	var output IsAuthorizedTaskSuccessOutput
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return false, err
	}
	return output.Authorized, nil
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
