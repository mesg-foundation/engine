package provider

import (
	"github.com/mesg-foundation/core/service/importer"
)

const (
	// publishVersion is the version used to publish the services to the marketplace
	publishVersion = 1

	// deploymentType is the type of deployment used for the service
	deploymentType = "ipfs"

	// marketplaceServiceKey is the key of the marketplace service
	marketplaceServiceKey = "marketplace"
)

// ErrorOutput is the output for any task that fails.
type ErrorOutput struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e ErrorOutput) Error() string {
	return e.Message
}

type ManifestData struct {
	Version    int                        `json:"version"`
	Definition importer.ServiceDefinition `json:"definition"`
	Readme     string                     `json:"readme,omitempty"`
	Service    struct {
		Deployment struct {
			Type   string `json:"type"`
			Source string `json:"source"`
		} `json:"deployment"`
	} `json:"service"`
}

// TransactionTaskInputs is the inputs for any task that create a transaction.
type TransactionTaskInputs struct {
	From string `json:"from"`
	// Gas      string `json:"gas"` // omitempty
	// GasPrice string `json:"gasPrice"` // omitempty
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

// PublishServiceVersionTaskInputs is the inputs of the task publish service version.
type PublishServiceVersionTaskInputs struct {
	*TransactionTaskInputs
	Sid              string `json:"sid"`
	VersionHash      string `json:"versionHash"`
	Manifest         string `json:"manifest"`
	ManifestProtocol string `json:"manifestProtocol"`
}

// CreateServiceOfferTaskInputs is the inputs of the task create service offer.
type CreateServiceOfferTaskInputs struct {
	*TransactionTaskInputs
	Sid      string `json:"sid"`
	Price    string `json:"price"`
	Duration string `json:"duration"`
}

// DisableServiceOfferTaskInputs is the inputs of the task create service offer.
type DisableServiceOfferTaskInputs struct {
	*TransactionTaskInputs
	Sid        string `json:"sid"`
	OfferIndex string `json:"offerIndex"`
}

// PurchaseTaskInputs is the inputs of the task purchase.
type PurchaseTaskInputs struct {
	*TransactionTaskInputs
	Sid        string `json:"sid"`
	OfferIndex string `json:"offerIndex"`
}

// PurchaseTaskOutputs is the output of the task purchase.
type PurchaseTaskOutputs struct {
	Transactions []*Transaction `json:"transactions"`
}

// TransferServiceOwnershipTaskInputs is the inputs of the task transfer service ownership.
type TransferServiceOwnershipTaskInputs struct {
	*TransactionTaskInputs
	Sid      string `json:"sid"`
	NewOwner string `json:"newOwner"`
}

// SendSignedTransactionTaskInputs is the inputs of the task send signed transaction.
type SendSignedTransactionTaskInputs struct {
	SignedTransaction string `json:"signedTransaction"`
}

// TransactionReceipt is the success output of task send signed transaction.
type TransactionReceipt struct {
	Receipt struct {
		BlockNumber      uint   `json:"blockNumber"`
		From             string `json:"from"`
		GasUsed          uint   `json:"gasUsed"`
		Status           bool   `json:"status"`
		To               string `json:"to"`
		TransactionHash  string `json:"transactionHash"`
		TransactionIndex uint   `json:"transactionIndex"`
	} `json:"receipt"`
}

// ServiceVersionInputs is the input for get service version task.
type ServiceVersionInputs struct {
	Sid         string `json:"sid"`
	VersionHash string `json:"versionHash"`
}

// ServiceVersionSuccessOutput is the input for get service version task.
type ServiceVersionSuccessOutput struct {
	VersionHash      string       `json:"versionHash"`
	Manifest         string       `json:"manifest"`
	ManifestProtocol string       `json:"manifestProtocol"`
	ManifestData     ManifestData `json:"manifestData"`
}

// GetServiceTaskInputs is the inputs of the task service exist.
type GetServiceTaskInputs struct {
	Sid string `json:"sid"`
}

// MarketplaceService is the success output of task service exist.
type MarketplaceService struct {
	Sid      string `json:"sid"`
	Owner    string `json:"owner"`
	Versions []struct {
		VersionHash      string       `json:"versionHash"`
		Manifest         string       `json:"manifest"`
		ManifestProtocol string       `json:"manifestProtocol"`
		ManifestData     ManifestData `json:"manifestData"`
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

// IsAuthorizedInputs is the inputs of the task check for deployment.
type IsAuthorizedInputs struct {
	Sid         string   `json:"sid"`
	VersionHash string   `json:"versionHash"`
	Addresses   []string `json:"addresses"`
}

// IsAuthorizedSuccessOutput is the success output of task check for deployment.
type IsAuthorizedSuccessOutput struct {
	Authorized bool   `json:"authorized"`
	Sid        string `json:"sid"`
	Source     string `json:"source"`
	Type       string `json:"type"`
}
