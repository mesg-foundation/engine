package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type marketplacePublishCmd struct {
	baseCmd

	noPassphrase bool
	passphrase   string
	address      string
	path         string
	manifest     provider.ManifestData

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
		Args:    cobra.MaximumNArgs(1),
	})
	c.cmd.Flags().StringVarP(&c.address, "address", "a", c.address, "Address to use to publish the service")
	c.cmd.Flags().BoolVarP(&c.noPassphrase, "no-passphrase", "n", c.noPassphrase, "Leave passphrase empty")
	c.cmd.Flags().StringVarP(&c.passphrase, "passphrase", "p", c.passphrase, "Passphrase to use with the account")
	return c
}

func (c *marketplacePublishCmd) preRunE(cmd *cobra.Command, args []string) error {
	var err error
	if c.address == "" {
		if err := askInput("Enter the address to use", &c.address); err != nil {
			return err
		}
	}
	if !c.noPassphrase && c.passphrase == "" {
		if err := askPass("Enter passphrase", &c.passphrase); err != nil {
			return err
		}
	}
	c.path = getFirstOrDefault(args)
	c.manifest, err = c.e.CreateManifest(c.path)
	if err != nil {
		return err
	}

	var confirmed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Are you sure to publish service %q from path %q using address %q?", c.manifest.Definition.Sid, c.path, c.address),
	}, &confirmed, nil); err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("cancel")
	}

	return nil
}

func (c *marketplacePublishCmd) runE(cmd *cobra.Command, args []string) error {
	// TODO: check if service already exist on the marketplace
	fmt.Println("CreateService transaction")
	tx, err := c.e.CreateService(c.manifest.Definition.Sid, c.address)
	if err != nil {
		return err
	}
	fmt.Println("CreateService transaction done", tx)

	fmt.Println("Sign transaction")
	// TODO: facto sign and SendSignedTransaction functions
	signedTransaction, err := c.e.Sign(c.address, c.passphrase, tx)
	if err != nil {
		return err
	}
	fmt.Println("Transaction signed", signedTransaction)

	fmt.Println("Send transaction")
	receipt, err := c.e.SendSignedTransaction(signedTransaction)
	if err != nil {
		return err
	}
	fmt.Println("Transaction sent", receipt.Receipt.TransactionHash)

	fmt.Println("Upload service source")
	// TODO: add a progress for the upload
	manifestProtocol, manifestSource, err := c.e.UploadServiceFiles(c.path, c.manifest)
	fmt.Println("Upload done. https://gateway.ipfs.io/ipfs/" + manifestSource)

	fmt.Println("CreateServiceVersion transaction")
	tx, err = c.e.CreateServiceVersion(c.manifest.Definition.Sid, "0x0001", manifestSource, manifestProtocol, c.address)
	if err != nil {
		return err
	}
	fmt.Println("CreateServiceVersion transaction done", tx)

	fmt.Println("Sign transaction")
	signedTransaction, err = c.e.Sign(c.address, c.passphrase, tx)
	if err != nil {
		return err
	}
	fmt.Println("Transaction signed", signedTransaction)

	fmt.Println("Send transaction")
	receipt, err = c.e.SendSignedTransaction(signedTransaction)
	if err != nil {
		return err
	}
	fmt.Println("Transaction sent", receipt.Receipt.TransactionHash)

	fmt.Println("done!")

	return nil
}
