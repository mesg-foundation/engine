package main

import (
	"bufio"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func signCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign <name> <message>",
		Short: "Sign the given message",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("must have exactly 2 args: name of the account to use and the message to sign")
			}
			buf := bufio.NewReader(cmd.InOrStdin())
			kb, err := keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), buf)
			if err != nil {
				return err
			}
			signature, _, err := kb.Sign(args[0], "", []byte(args[1]))
			if err != nil {
				return err
			}
			fmt.Println(hex.EncodeToString(signature))
			return nil
		},
	}
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	viper.BindPFlag(flags.FlagKeyringBackend, cmd.Flags().Lookup(flags.FlagKeyringBackend))
	return cmd
}

func verifySignCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verifySign <name> <message>",
		Short: "Verify the signature of the given message",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("must have exactly 2 args: name of the account to use and the message used to sign, and the signature")
			}
			buf := bufio.NewReader(cmd.InOrStdin())
			kb, err := keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), buf)
			if err != nil {
				return err
			}
			acc, err := kb.Get(args[0])
			if err != nil {
				return err
			}
			signature, err := hex.DecodeString(args[2])
			if err != nil {
				return err
			}
			verify := acc.GetPubKey().VerifyBytes([]byte(args[1]), signature)
			if verify {
				fmt.Println("verification succeed")
			} else {
				fmt.Println("verification failed")
			}
			return nil
		},
	}
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	viper.BindPFlag(flags.FlagKeyringBackend, cmd.Flags().Lookup(flags.FlagKeyringBackend))
	return cmd
}
