package commands

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type marketplacePublishCmd struct {
	baseCmd

	noPassphrase bool
	passphrase   string
	address      string
	path         string
	manifest     provider.ManifestData
	service      *provider.MarketplaceService

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
	var (
		err      error
		question = "a new service"
	)
	if c.address == "" {
		if err := askInput("Enter the address to use", &c.address); err != nil {
			return err
		}
	}
	c.address = strings.ToLower(c.address)
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

	var exist bool
	pretty.Progress("Checking if service exist...", func() {
		exist, err = c.e.ServiceExist(c.manifest.Definition.Sid)
	})
	if err != nil {
		return err
	}
	if exist {
		question = "a new version of service"
		pretty.Progress("Getting service data...", func() {
			c.service, err = c.e.GetService(c.manifest.Definition.Sid)
		})
		if err != nil {
			return err
		}
		if c.service.Owner != c.address {
			return fmt.Errorf("the service's owner (%q) is different than the specified address", c.service.Owner)
		}
	}

	var confirmed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Are you sure to publish %s %q from path %q using address %q?", question, c.manifest.Definition.Sid, c.path, c.address),
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
		hash             string
		manifestProtocol string
		manifestSource   string
	)
	if c.service == nil {
		pretty.Progress("Publishing service on the marketplace...", func() {
			tx, err = c.e.CreateService(c.manifest.Definition.Sid, c.address)
			if err != nil {
				return
			}
			_, err = c.signAndSendTransaction(tx)
		})
		if err != nil {
			return err
		}
		fmt.Printf("%s Service created with success\n", pretty.SuccessSign)
	}

	pretty.Progress("Uploading new version source code on the marketplace...", func() {
		// TODO: add a progress for the upload
		manifestProtocol, manifestSource, err = c.e.UploadServiceFiles(c.path, c.manifest)
	})
	if err != nil {
		return err
	}
	pretty.Progress("Publishing new version on the marketplace...", func() {
		h := sha256.New()
		h.Write([]byte(manifestSource))
		hash = fmt.Sprintf("0x%x", h.Sum(nil))
		tx, err = c.e.CreateServiceVersion(c.manifest.Definition.Sid, hash, manifestSource, manifestProtocol, c.address)
		if err != nil {
			return
		}
		_, err = c.signAndSendTransaction(tx)
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Version created with success\n", pretty.SuccessSign)
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(c.manifest.Definition.Sid))
	sidHash := fmt.Sprintf("0x%x", h.Sum(nil))
	fmt.Printf("%s See it on the marketplace:\nhttps://marketplace.mesg.com/services/%s/%s\n", pretty.SuccessSign, sidHash, hash)

	return nil
}

func (c *marketplacePublishCmd) signAndSendTransaction(transaction *provider.Transaction) (*provider.TransactionReceipt, error) {
	signedTransaction, err := c.e.Sign(c.address, c.passphrase, transaction)
	if err != nil {
		return nil, err
	}
	receipt, err := c.e.SendSignedTransaction(signedTransaction)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
