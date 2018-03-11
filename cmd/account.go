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
	"github.com/mesg-foundation/application/cmd/account"
	"github.com/spf13/cobra"
)

// Account is the command to manage all account activity
var Account = &cobra.Command{
	Use:   "account",
	Short: "Manage your MESG accounts",
}

func init() {
	Account.AddCommand(cmdAccount.Create)
	Account.AddCommand(cmdAccount.List)
	Account.AddCommand(cmdAccount.Delete)

	RootCmd.AddCommand(Account)
}
