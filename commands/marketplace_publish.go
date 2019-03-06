package commands

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type marketplacePublishCmd struct {
	baseMarketplaceCmd

	path     string
	manifest provider.ManifestData
	service  *provider.MarketplaceService

	e Executor
}

func newMarketplacePublishCmd(e Executor) *marketplacePublishCmd {
	c := &marketplacePublishCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "publish",
		Short:   "Publish a service on the MESG Marketplace",
		Example: `mesg-core marketplace publish PATH_TO_SERVICE`,
		PreRunE: c.preRunE,
		RunE:    c.runE,
		Args:    cobra.MinimumNArgs(1),
	})
	c.setupFlags()
	return c
}

func (c *marketplacePublishCmd) preRunE(cmd *cobra.Command, args []string) error {
	var (
		err      error
		question = "a new service"
	)

	if err := c.askAccountAndPassphrase(); err != nil {
		return err
	}

	c.path = getFirstOrDefault(args)
	c.manifest, err = c.e.CreateManifest(c.path)
	if err != nil {
		return err
	}

	pretty.Progress("Getting service data...", func() {
		c.service, err = c.e.GetService(c.manifest.Definition.Sid)
	})
	outputError, isOutputError := err.(provider.ErrorOutput)
	if isOutputError {
		err = nil
		if outputError.Code != "notFound" {
			question = "a new version of service"
			if !strings.EqualFold(c.service.Owner, c.account) {
				return fmt.Errorf("the service's owner (%q) is different than the specified account", c.service.Owner)
			}
		}
	}
	if err != nil {
		return err
	}

	var confirmed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Are you sure to publish %s %q from path %q using account %q?", question, c.manifest.Definition.Sid, c.path, c.account),
	}, &confirmed, nil); err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("cancelled")
	}

	return nil
}

func (c *marketplacePublishCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		tx               *provider.Transaction
		err              error
		versionHash      string
		manifestProtocol string
		manifestSource   string
	)

	// TODO: the hash of the version should be the hash after deployment

	pretty.Progress("Uploading service source code on the marketplace...", func() {
		// TODO: add a progress for the upload
		manifestProtocol, manifestSource, err = c.e.UploadServiceFiles(c.path, c.manifest)
	})
	if err != nil {
		return err
	}
	pretty.Progress("Publishing service on the marketplace...", func() {
		h := sha256.New()
		h.Write([]byte(manifestSource))
		versionHash = fmt.Sprintf("0x%x", h.Sum(nil))
		tx, err = c.e.PublishServiceVersion(c.manifest.Definition.Sid, versionHash, manifestSource, manifestProtocol, c.account)
		if err != nil {
			return
		}
		_, err = c.signAndSendTransaction(c.e, tx)
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Service published with success\n", pretty.SuccessSign)
	fmt.Printf("%s See it on the marketplace: https://marketplace.mesg.com/services/%s/%s\n", pretty.SuccessSign, c.manifest.Definition.Sid, versionHash)

	return nil
}
