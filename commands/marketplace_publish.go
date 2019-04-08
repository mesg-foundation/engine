package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/utils/readme"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

const (
	// marketplaceServiceHashVersion is the version of the service hash used by the core.
	marketplaceServiceHashVersion = "1"
)

type marketplacePublishCmd struct {
	baseMarketplaceCmd

	path    string
	service *coreapi.Service

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
	if err := c.askAccountAndPassphrase(); err != nil {
		return err
	}

	c.path = getFirstOrCurrentPath(args)

	var err error
	_, hash, err := deployService(c.e, c.path, map[string]string{})
	if err != nil {
		return err
	}
	c.service, err = c.e.ServiceByID(hash)
	if err != nil {
		return err
	}
	fmt.Printf("%s Service deployed with sid %s and hash %s\n", pretty.SuccessSign, pretty.Success(c.service.Sid), pretty.Success(c.service.Hash))

	var confirmed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Are you sure to publish a new version of service %q from path %q using account %q?", c.service.Sid, c.path, c.account),
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
		tx         provider.Transaction
		err        error
		deployment provider.MarketplaceDeployedSource
	)

	pretty.Progress("Uploading service source code on the marketplace...", func() {
		// TODO: add a progress for the upload
		deployment, err = c.e.UploadSource(c.path)
	})
	if err != nil {
		return err
	}
	readme, err := readme.Lookup(c.path)
	if err != nil {
		return err
	}
	pretty.Progress("Publishing service on the marketplace...", func() {
		tx, err = c.e.PublishServiceVersion(provider.MarketplaceManifestServiceData{
			Definition:  c.service,
			Hash:        c.service.Hash,
			HashVersion: marketplaceServiceHashVersion,
			Readme:      readme,
			Deployment:  deployment,
		}, c.account)
		if err != nil {
			return
		}
		_, err = c.signAndSendTransaction(c.e, tx)
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Service published with success\n", pretty.SuccessSign)
	fmt.Printf("%s See it on the marketplace: https://marketplace.mesg.com/services/%s\n", pretty.SuccessSign, c.service.Sid)

	fmt.Printf("%s To create a service offer, execute the command:\n\tmesg-core marketplace create-offer %s\n", pretty.SuccessSign, c.service.Sid)

	return nil
}
