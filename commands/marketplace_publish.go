package commands

import (
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
	manifest provider.MarketplaceManifestData
	service  provider.MarketplaceService

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

	c.path = getFirstOrCurrentPath(args)

	sid, hash, err := deployService(c.e, c.path, map[string]string{})
	if err != nil {
		return err
	}
	fmt.Printf("%s Service deployed with sid %s and hash %s\n", pretty.SuccessSign, pretty.Success(sid), pretty.Success(hash))

	c.manifest, err = c.e.CreateManifest(c.path, hash)
	if err != nil {
		return err
	}

	pretty.Progress("Getting service data...", func() {
		c.service, err = c.e.GetService(c.manifest.Service.Definition.Sid)
	})
	outputError, isMarketplaceError := err.(provider.MarketplaceErrorOutput)
	if err != nil && !isMarketplaceError {
		return err
	}
	if outputError.Code != "notFound" {
		question = "a new version of service"
		if !strings.EqualFold(c.service.Owner, c.account) {
			return fmt.Errorf("the service's owner %q is different than the specified account", c.service.Owner)
		}
	}

	var confirmed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Are you sure to publish %s %q from path %q using account %q?", question, c.manifest.Service.Definition.Sid, c.path, c.account),
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
		tx               provider.Transaction
		err              error
		manifestProtocol string
		manifestSource   string
	)

	pretty.Progress("Uploading service source code on the marketplace...", func() {
		// TODO: add a progress for the upload
		manifestProtocol, manifestSource, err = c.e.UploadServiceFiles(c.path, c.manifest)
	})
	if err != nil {
		return err
	}
	pretty.Progress("Publishing service on the marketplace...", func() {
		tx, err = c.e.PublishServiceVersion(c.manifest.Service.Definition.Sid, manifestSource, manifestProtocol, c.account)
		if err != nil {
			return
		}
		_, err = c.signAndSendTransaction(c.e, tx)
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Service published with success\n", pretty.SuccessSign)
	fmt.Printf("%s See it on the marketplace: https://marketplace.mesg.com/services/%s\n", pretty.SuccessSign, c.manifest.Service.Definition.Sid)

	fmt.Printf("%s To create a service offer, execute the command:\n\tmesg-core marketplace create-offer %s\n", pretty.SuccessSign, c.manifest.Service.Definition.Sid)

	return nil
}
