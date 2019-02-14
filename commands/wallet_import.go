package commands

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/x/xjson"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type walletImportCmd struct {
	baseWalletCmd

	importType string
	privateKey string
	jsonFile   string
	account    provider.EncryptedKeyJSONV3

	e WalletExecutor
}

func newWalletImportCmd(e WalletExecutor) *walletImportCmd {
	c := &walletImportCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "import ",
		Short:   "import an account",
		Long:    "import an account to wallet",
		Example: "mesg-core wallet import",
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})

	c.setupFlags()
	c.cmd.Flags().StringVarP(&c.privateKey, "private-key", "k", c.privateKey, "Private key to import")
	c.cmd.Flags().StringVarP(&c.jsonFile, "json", "j", c.jsonFile, "Filepath to the JSON file to import")
	return c
}

func (c *walletImportCmd) preRunE(cmd *cobra.Command, args []string) error {
	if c.jsonFile == "" && c.privateKey == "" {
		survey.AskOne(&survey.Select{
			Message: "How to import the account:",
			Options: []string{"JSON file", "private key"},
		}, &c.importType, nil)
		if c.importType == "private key" {
			if err := survey.AskOne(&survey.Password{
				Message: "Enter the private key to import",
			}, &c.privateKey, survey.MinLength(1)); err != nil {
				return err
			}
		} else {
			if err := survey.AskOne(&survey.Input{
				Message: "Enter the path to the JSON file to import",
			}, &c.jsonFile, survey.MinLength(1)); err != nil {
				return err
			}
		}
	}
	if c.jsonFile != "" {
		content, err := xjson.ReadFile(c.jsonFile)
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(content), &c.account); err != nil {
			return err
		}
	}
	if !c.noPassphrase && c.passphrase == "" {
		if err := survey.AskOne(&survey.Password{
			Message: "Enter passphrase",
		}, &c.passphrase, survey.MinLength(1)); err != nil {
			return err
		}
	}
	return nil
}

func (c *walletImportCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		address common.Address
		err     error
	)
	if c.privateKey != "" {
		address, err = c.e.ImportFromPrivateKey(c.privateKey, c.passphrase)
	} else {
		address, err = c.e.Import(c.account, c.passphrase)
	}
	if err != nil {
		return err
	}
	fmt.Printf("%s Wallet imported with address %s\n", pretty.SuccessSign, address.String())
	return nil
}
