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

// PublishServiceVersion executes the create service version task
func (p *MarketplaceProvider) PublishServiceVersion(sid, versionHash, manifest, manifestProtocol, from string) (*Transaction, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "publishServiceVersion", &PublishServiceVersionTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{From: from},
		Sid:                   sid,
		VersionHash:           versionHash,
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
func (p *MarketplaceProvider) CreateServiceOffer(sid string, price string, duration string, from string) (*Transaction, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "createServiceOffer", &CreateServiceOfferTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{From: from},
		Sid:                   sid,
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
func (p *MarketplaceProvider) DisableServiceOffer(sid, offerIndex, from string) (*Transaction, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "disableServiceOffer", &DisableServiceOfferTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{From: from},
		Sid:                   sid,
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
func (p *MarketplaceProvider) Purchase(sid, offerIndex, from string) (*PurchaseTaskOutputs, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "purchase", &PurchaseTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{From: from},
		Sid:                   sid,
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

	var output PurchaseTaskOutputs
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// TransferServiceOwnership executes the task transfer service ownership.
func (p *MarketplaceProvider) TransferServiceOwnership(sid, newOwner, from string) (*Transaction, error) {
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "transferServiceOwnership", &TransferServiceOwnershipTaskInputs{
		TransactionTaskInputs: &TransactionTaskInputs{From: from},
		Sid:                   sid,
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
func (p *MarketplaceProvider) SendSignedTransaction(signedTransaction string) (*TransactionReceipt, error) {
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

	var output TransactionReceipt
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// GetService executes the task get service.
func (p *MarketplaceProvider) GetService(sid string) (*MarketplaceService, error) {
	var output MarketplaceService
	r, err := p.client.ExecuteAndListen(MarketplaceServiceID, "getService", &GetServiceTaskInputs{
		Sid: sid,
	})
	if err != nil {
		return &output, err
	}

	if r.OutputKey == "error" {
		var outputError ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &outputError); err != nil {
			return &output, err
		}
		return &output, outputError
	}

	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return &output, err
	}
	return &output, nil
}

// UploadServiceFiles upload the tarball and the definition file, and returns the address of the definition file
func (p *MarketplaceProvider) UploadServiceFiles(path string, manifest ManifestData) (protocol string, source string, err error) {
	// upload service source to IPFS
	tar, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return "", "", err
	}
	tarballResponse, err := ipfs.Add("tarball", tar)
	if err != nil {
		return "", "", err
	}

	manifest.Service.Deployment.Type = DeploymentType
	manifest.Service.Deployment.Source = tarballResponse.Hash

	// upload manifest
	manifestData, err := json.Marshal(manifest)
	if err != nil {
		return "", "", err
	}
	definitionResponse, err := ipfs.Add("manifest", bytes.NewReader(manifestData))
	if err != nil {
		return "", "", err
	}

	return DeploymentType, definitionResponse.Hash, nil
}

func (p *MarketplaceProvider) CreateManifest(path string) (ManifestData, error) {
	var data ManifestData
	definition, err := importer.From(path)
	if err != nil {
		return data, err
	}
	data.Version = PublishVersion
	data.Definition = *definition
	data.Readme, err = readme.Lookup(path)
	if err != nil {
		return data, err
	}
	return data, nil
}
