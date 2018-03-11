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
)

func startCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "start",
		Short:   "Start a service",
		Args:    cobra.MinimumNArgs(1),
		Example: "mesg-cli service start --stake 100 --duration 10 ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			stake := cmd.Flag("stake").Value
			duration := cmd.Flag("duration").Value
			fmt.Println("service start called", args, stake, duration)
		},
	}
	cmd.Flags().Float64P("stake", "s", 0, "The number of MESG to put on stake")
	cmd.Flags().IntP("duration", "d", 0, "The amount of time you will be running this service for (in hours)")
	cmd.MarkFlagRequired("stake")
	cmd.MarkFlagRequired("duration")
	return cmd
}

func stopCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "stop",
		Short:   "Stop a service",
		Args:    cobra.MinimumNArgs(1),
		Example: "mesg-cli service stop ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("service stop called", args)
		},
	}
}

func pauseCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "pause",
		Short:   "Pause a service",
		Args:    cobra.MinimumNArgs(1),
		Example: "mesg-cli service pause ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("service pause called", args)
		},
	}
}

func resumeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "resume",
		Short:   "Resume a service",
		Args:    cobra.MinimumNArgs(1),
		Example: "mesg-cli service resume ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("service resume called", args)
		},
	}
}

func init() {
	var serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Manage the services you are running",
	}
	serviceCmd.AddCommand(startCmd())
	serviceCmd.AddCommand(stopCmd())
	serviceCmd.AddCommand(pauseCmd())
	serviceCmd.AddCommand(resumeCmd())
	RootCmd.AddCommand(serviceCmd)
}
