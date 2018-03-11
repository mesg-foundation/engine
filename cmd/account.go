// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func createAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Short:   "Create a new account",
		Example: "mesg-cli account create",
		Run: func(cmd *cobra.Command, args []string) {
			var password, passwordConfirmation, privateSeed string
			survey.AskOne(&survey.Password{Message: "Please set a password ?"}, &password, nil)
			survey.AskOne(&survey.Password{Message: "Repeat your password ?"}, &passwordConfirmation, nil)
			if password != passwordConfirmation {
				fmt.Println("Password confirmation invalid")
				return
			}
			fmt.Println("Generating secure key")
			survey.AskOne(&survey.Input{Message: "Repeat your private seed"}, &privateSeed, nil)
			// TODO add real account creation
			fmt.Println("creating account", password)
		},
	}
}

func listAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List all the accounts on this computer",
		Example: "mesg-cli account list",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO add real listing
			fmt.Println("0x0000000000000000000000000000000000000000")
			fmt.Println("0x0000000000000000000000000000000000000001")
		},
	}
}

func deleteAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Short:   "Delete an account",
		Example: "mesg-cli service delete 0x0000000000000000000000000000000000000000",
		Run: func(cmd *cobra.Command, args []string) {
			var confirm bool
			var account = ""
			if len(args) > 0 {
				account = args[0]
			}
			if account == "" {
				// TODO add real list
				accounts := []string{"0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000001"}
				survey.AskOne(&survey.Select{
					Message: "Choose the account you want to delete",
					Default: accounts[0],
					Options: accounts,
				}, &account, nil)
			}
			survey.AskOne(&survey.Confirm{Message: "Are you sure ? You can always re-import an account with your private seed"}, &confirm, nil)
			if confirm {
				// TODO add real deletion
				fmt.Println("account deleted", account)
			}
		},
	}
}

func init() {
	var accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Manage your MESG accounts",
	}
	accountCmd.AddCommand(createAccountCmd())
	accountCmd.AddCommand(listAccountCmd())
	accountCmd.AddCommand(deleteAccountCmd())
	// accountCmd.AddCommand(resumeCmd())
	RootCmd.AddCommand(accountCmd)
}
