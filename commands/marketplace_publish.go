package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type marketplacePublishCmd struct {
	baseCmd

	path string

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
	return c
}

func (c *marketplacePublishCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrDefault(args)
	return nil
}

func (c *marketplacePublishCmd) runE(cmd *cobra.Command, args []string) error {
	from := "0xf3C21FD07B1D4c40d3cE6EfaC81a3E49f6c04592"
	passphrase := "1"

	tx, err := c.e.CreateService("test1", from)
	if err != nil {
		return err
	}
	fmt.Println("tx", tx)

	signedTransaction, err := c.e.Sign(from, passphrase, tx)
	if err != nil {
		return err
	}
	fmt.Println("transaction signed", signedTransaction)

	receipt, err := c.e.SendSignedTransaction(signedTransaction)
	if err != nil {
		return err
	}
	fmt.Println("service created", receipt)

	// definition, err := c.e.PublishDefinitionFile(c.path)
	// fmt.Println("upload done. https://gateway.ipfs.io/ipfs/" + definition)

	// tx, err = c.e.CreateServiceVersion(sidHash, hash, definition, "ipfs", from)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("tx", tx)

	// fmt.Println("published")

	return nil
}
