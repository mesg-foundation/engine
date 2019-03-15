package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"
)

type rootMarketplaceCmd struct {
	baseCmd
}

func newRootMarketplaceCmd(e Executor) *rootMarketplaceCmd {
	c := &rootMarketplaceCmd{}
	c.cmd = newCommand(&cobra.Command{
		Use:   "marketplace",
		Short: "Interact with the Marketplace",
	})

	c.cmd.AddCommand(
		newMarketplacePublishCmd(e).cmd,
		newMarketplaceCreateOfferCmd(e).cmd,
		newMarketplacePurchaseCmd(e).cmd,
	)
	return c
}

// baseMarketplaceCmd is basic command for marketplace subcommands
// that require passphrase.
type baseMarketplaceCmd struct {
	baseCmd
	account      string
	noPassphrase bool
	passphrase   string
}

func (c *baseMarketplaceCmd) setupFlags() {
	c.cmd.Flags().StringVarP(&c.account, "account", "a", c.account, "Account to use")
	c.cmd.Flags().BoolVarP(&c.noPassphrase, "no-passphrase", "n", c.noPassphrase, "Leave passphrase empty")
	c.cmd.Flags().StringVarP(&c.passphrase, "passphrase", "p", c.passphrase, "Passphrase to decrypt the account")
}

func (c *baseMarketplaceCmd) askAccountAndPassphrase() error {
	if c.account == "" {
		if err := askInput("Enter the account to use", &c.account); err != nil {
			return err
		}
	}
	if !c.noPassphrase && c.passphrase == "" {
		if err := askPass("Enter passphrase", &c.passphrase); err != nil {
			return err
		}
	}
	return nil
}

func (c *baseMarketplaceCmd) signAndSendTransaction(e Executor, transaction provider.Transaction) (provider.TransactionReceipt, error) {
	signedTransaction, err := e.Sign(c.account, c.passphrase, transaction)
	if err != nil {
		return provider.TransactionReceipt{}, err
	}
	return e.SendSignedTransaction(signedTransaction)
}

func (c *baseMarketplaceCmd) sha3(s string) string {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(s))
	return fmt.Sprintf("0x%x", h.Sum(nil))
}
