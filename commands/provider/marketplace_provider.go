package provider

import (
	"bytes"
	"encoding/json"

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
func (p *MarketplaceProvider) PublishServiceVersion(sid, manifest, manifestProtocol, from string) (Transaction, error) {
	input := marketplacePublishServiceVersionTaskInputs{
		marketplaceTransactionTaskInputs: marketplaceTransactionTaskInputs{From: from},
		Sid:                              sid,
		Manifest:                         manifest,
		ManifestProtocol:                 manifestProtocol,
	}
	var output Transaction
	return output, p.call("publishServiceVersion", input, &output)
}

// CreateServiceOffer executes the create service offer task
func (p *MarketplaceProvider) CreateServiceOffer(sid string, price string, duration string, from string) (Transaction, error) {
	input := marketplaceCreateServiceOfferTaskInputs{
		marketplaceTransactionTaskInputs: marketplaceTransactionTaskInputs{From: from},
		Sid:                              sid,
		Price:                            price,
		Duration:                         duration,
	}
	var output Transaction
	return output, p.call("createServiceOffer", input, &output)
}

// DisableServiceOffer executes the disable service offer task
func (p *MarketplaceProvider) DisableServiceOffer(sid, offerIndex, from string) (Transaction, error) {
	input := marketplaceDisableServiceOfferTaskInputs{
		marketplaceTransactionTaskInputs: marketplaceTransactionTaskInputs{From: from},
		Sid:                              sid,
		OfferIndex:                       offerIndex,
	}
	var output Transaction
	return output, p.call("disableServiceOffer", input, &output)
}

// Purchase executes the purchase task
func (p *MarketplaceProvider) Purchase(sid, offerIndex, from string) ([]Transaction, error) {
	input := marketplacePurchaseTaskInputs{
		marketplaceTransactionTaskInputs: marketplaceTransactionTaskInputs{From: from},
		Sid:                              sid,
		OfferIndex:                       offerIndex,
	}
	var output marketplacePurchaseTaskOutputs
	return output.Transactions, p.call("purchase", input, &output)
}

// TransferServiceOwnership executes the task transfer service ownership.
func (p *MarketplaceProvider) TransferServiceOwnership(sid, newOwner, from string) (Transaction, error) {
	input := marketplaceTransferServiceOwnershipTaskInputs{
		marketplaceTransactionTaskInputs: marketplaceTransactionTaskInputs{From: from},
		Sid:                              sid,
		NewOwner:                         newOwner,
	}
	var output Transaction
	return output, p.call("transferServiceOwnership", input, &output)
}

// SendSignedTransaction executes the task send signed transaction.
func (p *MarketplaceProvider) SendSignedTransaction(signedTransaction string) (TransactionReceipt, error) {
	input := marketplaceSendSignedTransactionTaskInputs{
		SignedTransaction: signedTransaction,
	}
	var output TransactionReceipt
	return output, p.call("sendSignedTransaction", input, &output)
}

// GetService executes the task get service.
func (p *MarketplaceProvider) GetService(sid string) (MarketplaceService, error) {
	input := marketplaceGetServiceTaskInputs{
		Sid: sid,
	}
	var output MarketplaceService
	return output, p.call("getService", input, &output)
}

// IsAuthorized executes the task IsAuthorized.
func (p *MarketplaceProvider) IsAuthorized(sid string, versionHash string, addresses []string) (bool, string, string, string, error) {
	input := marketplaceIsAuthorizedInputs{
		VersionHash: versionHash,
		Sid:         sid,
		Addresses:   addresses,
	}
	var output marketplaceIsAuthorizedSuccessOutput
	return output.Authorized, output.Sid, output.Source, output.Type, p.call("isAuthorized", input, &output)
}

// UploadServiceFiles upload the tarball and the definition file, and returns the address of the definition file
func (p *MarketplaceProvider) UploadServiceFiles(path string, manifest MarketplaceManifestData) (protocol string, source string, err error) {
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

	manifest.Service.Deployment.Type = marketplaceDeploymentType
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

	return marketplaceDeploymentType, definitionResponse.Hash, nil
}

func (p *MarketplaceProvider) CreateManifest(path string, hash string) (MarketplaceManifestData, error) {
	var data MarketplaceManifestData
	definition, err := importer.From(path)
	if err != nil {
		return data, err
	}
	data.Version = marketplacePublishVersion
	data.Service.Hash = hash
	data.Service.HashVersion = "1" // hardcoded for now
	data.Service.Definition = *definition
	data.Service.Readme, err = readme.Lookup(path)
	return data, err
}

func (p *MarketplaceProvider) call(task string, inputs interface{}, output interface{}) error {
	serviceHash, err := p.client.GetServiceHash(marketplaceServiceKey)
	if err != nil {
		return err
	}
	r, err := p.client.ExecuteAndListen(serviceHash, task, &inputs)
	if err != nil {
		return err
	}
	return p.parseResult(r, &output)
}

func (p *MarketplaceProvider) parseResult(r *coreapi.ResultData, output interface{}) error {
	if r.OutputKey == "error" {
		var outputError MarketplaceErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &outputError); err != nil {
			return err
		}
		return outputError
	}
	return json.Unmarshal([]byte(r.OutputData), &output)
}
