package provider

import (
	"encoding/json"

	"github.com/mesg-foundation/core/service/importer"
)

const (
	// marketplaceDeploymentType is the type of deployment used for the service.
	marketplaceDeploymentType = "ipfs"

	// marketplaceServiceKey is the key of the marketplace service.
	marketplaceServiceKey = "marketplace"
)

// MarketplaceErrorOutput is the output for any task that fails.
type MarketplaceErrorOutput struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e MarketplaceErrorOutput) Error() string {
	return e.Message
}

// MarketplaceService is the success output of task service exist.
type MarketplaceService struct {
	Sid      string `json:"sid"`
	Owner    string `json:"owner"`
	Versions []struct {
		VersionHash      string                  `json:"versionHash"`
		Manifest         string                  `json:"manifest"`
		ManifestProtocol string                  `json:"manifestProtocol"`
		ManifestData     MarketplaceManifestData `json:"manifestData,omitempty"`
	} `json:"versions"`
	Offers []struct {
		OfferIndex string `json:"offerIndex"`
		Price      string `json:"price"`
		Duration   string `json:"duration"`
		Active     bool   `json:"active"`
	} `json:"offers"`
	Purchases []struct {
		Purchaser string `json:"purchaser"`
		Expire    string `json:"expire"`
	} `json:"purchases"`
}

// SourceDeployment is the information related to a deployment
type SourceDeployment struct {
	Type   string `json:"type"`
	Source string `json:"source"`
}

// MarketplaceServiceData is the data present to the manifest and sent to create a new service's version
type MarketplaceServiceData struct {
	Definition  importer.ServiceDefinition `json:"definition"`
	Readme      string                     `json:"readme,omitempty"`
	Hash        string                     `json:"hash"`
	HashVersion string                     `json:"hashVersion"`
	Deployment  SourceDeployment           `json:"deployment"`
}

// MarketplaceManifestData struct {
type MarketplaceManifestData struct {
	Version string                 `json:"version"`
	Service MarketplaceServiceData `json:"service"`
}

// UnmarshalJSON overrides the default one to allow parsing malformed manifest data without returning error to user.
func (d *MarketplaceManifestData) UnmarshalJSON(data []byte) error {
	// the following temporary type prevents recursive cycling call when unmarshalling
	type tempType MarketplaceManifestData
	if err := json.Unmarshal(data, (*tempType)(d)); err != nil {
		*d = MarketplaceManifestData{}
	}
	return nil
}

type marketplaceTransactionTaskInputs struct {
	From     string `json:"from"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
}

type marketplacePublishServiceVersionTaskInputs struct {
	marketplaceTransactionTaskInputs
	Service MarketplaceServiceData `json:"service"`
}

type marketplaceCreateServiceOfferTaskInputs struct {
	marketplaceTransactionTaskInputs
	Sid      string `json:"sid"`
	Price    string `json:"price"`
	Duration string `json:"duration"`
}

type marketplacePurchaseTaskInputs struct {
	marketplaceTransactionTaskInputs
	Sid        string `json:"sid"`
	OfferIndex string `json:"offerIndex"`
}

type marketplacePurchaseTaskOutputs struct {
	Transactions []Transaction `json:"transactions"`
}

type marketplaceSendSignedTransactionTaskInputs struct {
	SignedTransaction string `json:"signedTransaction"`
}

type marketplaceGetServiceTaskInputs struct {
	Sid string `json:"sid"`
}

type marketplaceIsAuthorizedInputs struct {
	Sid         string   `json:"sid"`
	VersionHash string   `json:"versionHash"`
	Addresses   []string `json:"addresses"`
}

type marketplaceIsAuthorizedSuccessOutput struct {
	Authorized bool   `json:"authorized"`
	Sid        string `json:"sid"`
	Source     string `json:"source"`
	Type       string `json:"type"`
}
