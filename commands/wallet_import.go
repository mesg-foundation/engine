package commands

import (
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/x/xjson"
	"github.com/spf13/cobra"
)

type walletImportCmd struct {
	baseWalletCmd

	privateKey string
	jsonFile   string
	account    provider.WalletEncryptedKeyJSONV3

	e WalletExecutor
}

func newWalletImportCmd(e WalletExecutor) *walletImportCmd {
	c := &walletImportCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "import",
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
		if err := c.askImportType(); err != nil {
			return err
		}
	}
	if c.jsonFile != "" {
		if err := c.readJSONFile(c.jsonFile); err != nil {
			return err
		}
	}
	return c.askPassphrase()
}

func (c *walletImportCmd) askImportType() error {
	importType := []string{
		"json file",
		"private key",
	}
	var selectedImportType string
	if err := askSelect("How to import the account:", importType, &selectedImportType); err != nil {
		return err
	}
	if selectedImportType == importType[1] {
		if err := askPass("Enter the private key to import", &c.privateKey); err != nil {
			return err
		}
	}
	if selectedImportType == importType[0] {
		if err := askInput("Enter the path to the json file to import", &c.jsonFile); err != nil {
			return err
		}
	}
	return nil
}

func (c *walletImportCmd) readJSONFile(jsonFile string) error {
	content, err := xjson.ReadFile(jsonFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, &c.account)
}

func (c *walletImportCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		address string
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
	fmt.Printf("%s Account imported with address %s\n", pretty.SuccessSign, address)
	return nil
}
