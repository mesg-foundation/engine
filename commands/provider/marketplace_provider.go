package provider

import (
	"encoding/json"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/ipfs"
	"github.com/mesg-foundation/core/protobuf/coreapi"
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
func (p *MarketplaceProvider) PublishServiceVersion(service MarketplaceServiceData, from string) (Transaction, error) {
	input := marketplacePublishServiceVersionTaskInputs{
		marketplaceTransactionTaskInputs: marketplaceTransactionTaskInputs{From: from},
		Service:                          service,
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

// UploadSources upload the tarball, and returns the address of the uploaded sources
func (p *MarketplaceProvider) UploadSources(path string) (SourceDeployment, error) {
	// upload service source to IPFS
	tar, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return SourceDeployment{}, err
	}
	tarballResponse, err := ipfs.Add("tarball", tar)
	if err != nil {
		return SourceDeployment{}, err
	}
	return SourceDeployment{
		Type:   marketplaceDeploymentType,
		Source: tarballResponse.Hash,
	}, nil
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
