package provider

import (
	"encoding/json"

	"github.com/mesg-foundation/core/service/importer"
)

const (
	// marketplacePublishVersion is the version used to publish the services to the marketplace
	marketplacePublishVersion = "1"

	// marketplaceDeploymentType is the type of deployment used for the service
	marketplaceDeploymentType = "ipfs"

	// marketplaceServiceKey is the key of the marketplace service
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

// MarketplaceManifestData struct {
type MarketplaceManifestData struct {
	Version string `json:"version"`
	Service struct {
		Definition  importer.ServiceDefinition `json:"definition"`
		Readme      string                     `json:"readme,omitempty"`
		Hash        string                     `json:"hash"`
		HashVersion string                     `json:"hashVersion"`
		Deployment  struct {
			Type   string `json:"type"`
			Source string `json:"source"`
		} `json:"deployment"`
	} `json:"service"`
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
	Sid              string `json:"sid"`
	Manifest         string `json:"manifest"`
	ManifestProtocol string `json:"manifestProtocol"`
}

type marketplaceCreateServiceOfferTaskInputs struct {
	marketplaceTransactionTaskInputs
	Sid      string `json:"sid"`
	Price    string `json:"price"`
	Duration string `json:"duration"`
}

type marketplaceDisableServiceOfferTaskInputs struct {
	marketplaceTransactionTaskInputs
	Sid        string `json:"sid"`
	OfferIndex string `json:"offerIndex"`
}

type marketplacePurchaseTaskInputs struct {
	marketplaceTransactionTaskInputs
	Sid        string `json:"sid"`
	OfferIndex string `json:"offerIndex"`
}

type marketplacePurchaseTaskOutputs struct {
	Transactions []Transaction `json:"transactions"`
}

type marketplaceTransferServiceOwnershipTaskInputs struct {
	marketplaceTransactionTaskInputs
	Sid      string `json:"sid"`
	NewOwner string `json:"newOwner"`
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
