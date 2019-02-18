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

// IsAuthorizedTaskInputs is the inputs of the task is authorized.
type IsAuthorizedTaskInputs struct {
	Sid string `json:"sid"`
}

// IsAuthorizedTaskSuccessOutput is the success output of task authorized.
type IsAuthorizedTaskSuccessOutput struct {
	Authorized bool `json:"authorized"`
}

// ServiceVersionInputs is the input for get service version task.
type ServiceVersionInputs struct {
	Sid  string `json:"sid"`
	Hash string `json:"hash"`
}

// ServiceVersionSuccessOutput is the input for get service version task.
type ServiceVersionSuccessOutput struct {
	Hash             string       `json:"hash"`
	ManifestSource   string       `json:"manifestSource"`
	ManifestProtocol string       `json:"manifestProtocol"`
	Manifest         ManifestData `json:"manifest"`
}

// ServiceExistTaskInputs is the inputs of the task service exist.
type ServiceExistTaskInputs struct {
	Sid string `json:"sid"`
}

// ServiceExistTaskSuccessOutput is the success output of task service exist.
type ServiceExistTaskSuccessOutput struct {
	Exist bool `json:"exist"`
}

// GetServiceTaskInputs is the inputs of the task service exist.
type GetServiceTaskInputs struct {
	Sid string `json:"sid"`
}

// MarketplaceService is the success output of task service exist.
type MarketplaceService struct {
	Sid      string `json:"sid"`
	SidHash  string `json:"sidHash"`
	Owner    string `json:"owner"`
	Versions []struct {
		Hash             string       `json:"hash"`
		ManifestSource   string       `json:"manifestSource"`
		ManifestProtocol string       `json:"manifestProtocol"`
		Manifest         ManifestData `json:"manifest"`
	} `json:"versions"`
	Offers []struct {
		Index    string `json:"index"`
		Price    string `json:"price"`
		Duration string `json:"duration"`
		Active   bool   `json:"active"`
	} `json:"offers"`
	Purchases []struct {
		Purchaser string `json:"purchaser"`
		Expire    string `json:"expire"`
	} `json:"purchases"`
}

// CheckForDeploymentInputs is the inputs of the task check for deployment.
type CheckForDeploymentInputs struct {
	Hash      string   `json:"hash"`
	Addresses []string `json:"addresses"`
}

// CheckForDeploymentSuccessOutput is the success output of task check for deployment.
type CheckForDeploymentSuccessOutput struct {
	Authorized bool   `json:"authorized"`
	Source     string `json:"source"`
	Type       string `json:"type"`
}
