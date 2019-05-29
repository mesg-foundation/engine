package provider

import (
	"encoding/json"
	"time"

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

// PreparePublishServiceVersion executes the create service version task
func (p *MarketplaceProvider) PreparePublishServiceVersion(service MarketplaceManifestServiceData, from string) (Transaction, error) {
	input := marketplacePreparePublishServiceVersionTaskInputs{
		marketplacePrepareTaskInputs: marketplacePrepareTaskInputs{From: from},
		Service:                      service,
	}
	var output Transaction
	return output, p.call("preparePublishServiceVersion", input, &output)
}

// PublishPublishServiceVersion executes the task publish service version task
func (p *MarketplaceProvider) PublishPublishServiceVersion(signedTransaction string) (sid, versionHash, manifest, manifestProtocol string, err error) {
	input := marketplacePublishTaskInputs{
		SignedTransaction: signedTransaction,
	}
	var o marketplacePublishPublishServiceVersionTaskOutputs
	err = p.call("publishPublishServiceVersion", input, &o)
	return o.Sid, o.VersionHash, o.Manifest, o.ManifestProtocol, err
}

// PrepareCreateServiceOffer executes the create service offer task
func (p *MarketplaceProvider) PrepareCreateServiceOffer(sid string, price string, duration string, from string) (Transaction, error) {
	input := marketplacePrepareCreateServiceOfferTaskInputs{
		marketplacePrepareTaskInputs: marketplacePrepareTaskInputs{From: from},
		Sid:                          sid,
		Price:                        price,
		Duration:                     duration,
	}
	var output Transaction
	return output, p.call("prepareCreateServiceOffer", input, &output)
}

// PublishCreateServiceOffer executes the task publish service offer task
func (p *MarketplaceProvider) PublishCreateServiceOffer(signedTransaction string) (sid, offerIndex, price, duration string, err error) {
	input := marketplacePublishTaskInputs{
		SignedTransaction: signedTransaction,
	}
	var o marketplacePublishCreateServiceOfferTaskOutputs
	err = p.call("publishCreateServiceOffer", input, &o)
	return o.Sid, o.OfferIndex, o.Price, o.Duration, err
}

// PreparePurchase executes the purchase task
func (p *MarketplaceProvider) PreparePurchase(sid, offerIndex, from string) ([]Transaction, error) {
	input := marketplacePreparePurchaseTaskInputs{
		marketplacePrepareTaskInputs: marketplacePrepareTaskInputs{From: from},
		Sid:                          sid,
		OfferIndex:                   offerIndex,
	}
	var output marketplacePreparePurchaseTaskOutputs
	return output.Transactions, p.call("preparePurchase", input, &output)
}

// PublishPurchase executes the task publish service version task
func (p *MarketplaceProvider) PublishPurchase(signedTransactions []string) (sid, offerIndex, purchaser, price, duration string, expire time.Time, err error) {
	input := marketplacePublishPurchaseTaskInputs{
		SignedTransactions: signedTransactions,
	}
	var o marketplacePublishPurchaseTaskOutputs
	err = p.call("publishPurchase", input, &o)
	return o.Sid, o.OfferIndex, o.Purchaser, o.Price, o.Duration, o.Expire, err
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

// UploadSource upload the tarball, and returns the address of the uploaded sources
func (p *MarketplaceProvider) UploadSource(path string) (MarketplaceDeployedSource, error) {
	// upload service source to IPFS
	tar, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return MarketplaceDeployedSource{}, err
	}
	tarballResponse, err := ipfs.Add("tarball", tar)
	if err != nil {
		return MarketplaceDeployedSource{}, err
	}
	return MarketplaceDeployedSource{
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
	if r.Error != "" {
		var outputError MarketplaceErrorOutput
		if err := json.Unmarshal([]byte(r.Error), &outputError); err != nil {
			return err
		}
		return outputError
	}
	return json.Unmarshal([]byte(r.OutputData), &output)
}
