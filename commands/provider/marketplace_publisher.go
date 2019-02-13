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
	client coreapi.CoreClient
}

// NewMarketplaceProvider creates new MarketplaceProvider.
func NewMarketplaceProvider(c coreapi.CoreClient) *MarketplaceProvider {
	return &MarketplaceProvider{client: c}
}

// PublishVersion is the version used to publish the services to the marketplace
const PublishVersion = 1

// DeploymentType is the type of deployment used for the service
const DeploymentType = "ipfs"

type marketplaceData struct {
	Version int `json:"version"`
	Service struct {
		Deployment struct {
			Type   string `json:"type"`
			Source string `json:"source"`
		} `json:"deployment"`
	} `json:"service"`
	Definition *importer.ServiceDefinition `json:"definition"`
	Readme     string                      `json:"readme,omitempty"`
}

// MarketplacePublishDefinitionFile upload and publish the tarball and definition file and returns the address of the definition file
func (p *MarketplaceProvider) MarketplacePublishDefinitionFile(path string) (string, error) {
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
	var data marketplaceData
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
