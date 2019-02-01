package commands

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/ipfs"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/utils/readme"
	"github.com/spf13/cobra"
)

// PublishVersion is the version used to publish the services to the marketplace
const PublishVersion = 1

type servicePublishCmd struct {
	baseCmd

	path string
	ipfs *ipfs.IPFS

	e ServiceExecutor
}

type marketplaceData struct {
	Version int `json:"version"`
	Service struct {
		Deployment struct {
			Source string `json:"source"`
		} `json:"deployment"`
	} `json:"service"`
	Definition *importer.ServiceDefinition `json:"definition"`
	Readme     string                      `json:"readme,omitempty"`
}

func newServicePublishCmd(e ServiceExecutor) *servicePublishCmd {
	c := &servicePublishCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "publish",
		Short:   "Publish a service on the MESG Marketplace",
		Example: `mesg-core service publish PATH_TO_SERVICE`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
		Args:    cobra.MaximumNArgs(1),
		Hidden:  true, // TODO: Remove when this feature is finished
	})
	return c
}

func (c *servicePublishCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrDefault(args)
	c.ipfs = ipfs.New()
	return nil
}

func (c *servicePublishCmd) runE(cmd *cobra.Command, args []string) error {
	tar, err := archive.TarWithOptions(c.path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return err
	}
	tarballResponse, err := c.ipfs.Add("tarball", tar)
	if err != nil {
		return err
	}

	definitionFile, err := c.createDefinitionFile(tarballResponse.Hash)
	if err != nil {
		return err
	}

	definitionResponse, err := c.ipfs.Add("definition", bytes.NewReader(definitionFile))
	if err != nil {
		return err
	}

	fmt.Println("https://gateway.ipfs.io/ipfs/" + definitionResponse.Hash)

	return nil
}

func (c *servicePublishCmd) createDefinitionFile(tarballHash string) ([]byte, error) {
	definition, err := importer.From(c.path)
	if err != nil {
		return nil, err
	}
	var data marketplaceData
	data.Version = PublishVersion
	data.Service.Deployment.Source = tarballHash
	data.Readme, err = readme.LookupReadme(c.path)
	if err != nil {
		return nil, err
	}
	data.Definition = definition
	return json.Marshal(data)
}
