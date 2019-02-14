package provider

import (
	"github.com/mesg-foundation/core/service/importer"
)

const (
	// PublishVersion is the version used to publish the services to the marketplace
	PublishVersion = 1

	// DeploymentType is the type of deployment used for the service
	DeploymentType = "ipfs"

	// MarketplaceServiceID is the sid of the marketplace service
	MarketplaceServiceID = "marketplace"
)

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
	Sid              string `json:"sid"`
	Hash             string `json:"hash"`
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

// SendSignedTransactionTaskSuccessOutput is the success output of task send signed transaction.
type SendSignedTransactionTaskSuccessOutput struct {
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

// IsAuthorizedTaskInputs is the inputs of the task is authorized.
type IsAuthorizedTaskInputs struct {
	Sid string `json:"sid"`
}

// IsAuthorizedTaskSuccessOutput is the success output of task authorized.
type IsAuthorizedTaskSuccessOutput struct {
	Authorized bool `json:"authorized"`
}
