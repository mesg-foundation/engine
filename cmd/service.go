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

func confirm(cmd *cobra.Command, message string) bool {
	confirm := cmd.Flag("confirm").Value.String() == "true"
	if !confirm {
		survey.AskOne(&survey.Confirm{Message: message}, &confirm, nil)
	}
	return confirm
}

func startServiceCmd() *cobra.Command {
	var stake float64
	var duration int
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start a service",
		Args:    cobra.MinimumNArgs(1),
		Example: "mesg-cli service start --stake 100 --duration 10 ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			if stake == 0 {
				survey.AskOne(&survey.Input{
					Message: "How much do you want to stake (MESG) ?",
					Help:    "More details on the stake here",
					Default: "0",
				}, &stake, nil)
			}
			if duration == 0 {
				survey.AskOne(&survey.Input{
					Message: "How long will you run this service (hours) ?",
					Help:    "More details on the duration here",
					Default: "0",
				}, &duration, nil)
			}
			if !confirm(cmd, "Are you sure to run this service and stake your tokens ?") {
				return
			}
			// TODO stake && start service
			fmt.Println("service start called", args, stake, duration)
		},
	}
	cmd.Flags().BoolP("confirm", "c", false, "Confirm")
	cmd.Flags().Float64VarP(&stake, "stake", "s", 0, "The number of MESG to put on stake")
	cmd.Flags().IntVarP(&duration, "duration", "d", 0, "The amount of time you will be running this/those service(s) for (in hours)")
	return cmd
}

func stopServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a service",
		Long: `By stoping a service, your node will not process any other actions from this service.
/!\ This action will slash your stake if you didn't respect the duration`,
		Args:    cobra.MinimumNArgs(1),
		Example: "mesg-cli service stop ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			if !confirm(cmd, "Are you sure ? Your stake may be slashed !") {
				return
			}
			// TODO take stake && stop service
			// Is it really usefull to take the stake, the node will be offline anyway and we cannot trust the client
			fmt.Println("service stop called", args)
		},
	}
	cmd.Flags().BoolP("confirm", "c", false, "Confirm")
	return cmd
}

func pauseServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause",
		Short:   "Pause a service",
		Args:    cobra.MinimumNArgs(1),
		Example: "mesg-cli service pause ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			if !confirm(cmd, "Are you sure ?") {
				return
			}
			// TODO pause (onchain) and then stop the service
			fmt.Println("service pause called", args)
		},
	}
	cmd.Flags().BoolP("confirm", "c", false, "Confirm")
	return cmd
}

func resumeServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resume",
		Short:   "Resume a service",
		Args:    cobra.MinimumNArgs(1),
		Example: "mesg-cli service resume ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			if !confirm(cmd, "Are you sure ?") {
				return
			}
			// TODO start and when ready resume (onchain) the service
			fmt.Println("service resume called", args)
		},
	}
	cmd.Flags().BoolP("confirm", "c", false, "Confirm")
	return cmd
}

func init() {
	var serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Manage the services you are running",
	}
	serviceCmd.AddCommand(startServiceCmd())
	serviceCmd.AddCommand(stopServiceCmd())
	serviceCmd.AddCommand(pauseServiceCmd())
	serviceCmd.AddCommand(resumeServiceCmd())
	RootCmd.AddCommand(serviceCmd)
}
