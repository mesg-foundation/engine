package main

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/ext/xsignal"
	"github.com/mesg-foundation/engine/orchestrator"
	"github.com/mesg-foundation/engine/server/grpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	accName = "orchestrator"
	accPass = "password"
)

func orchestratorCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orchestrator",
		Short: "Orchestrator subcommands",
	}
	cmd.AddCommand(flags.GetCommands(
		startOrchestratorCmd(cdc),
	)...)
	return cmd
}

func startOrchestratorCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the Orchestrator",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			if viper.GetString(flagMnemonic) == "" {
				return fmt.Errorf("mnemonic is required. use flag --mnemonic or config file")
			}
			if cliCtx.ChainID == "" {
				return fmt.Errorf("chain-id is required. use flag --chain-id or config file")
			}

			logger := log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))
			client, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			// init rpc client
			logger.Info("Starting rpc client")
			if err := client.Start(); err != nil {
				return err
			}
			defer func() {
				logger.Info("Stopping rpc client")
				if err := client.Stop(); err != nil {
					logger.Error(err.Error())
				}
			}()

			kb := cosmos.NewInMemoryKeybase()
			if _, err := kb.CreateAccount(accName, viper.GetString(flagMnemonic), "", accPass, keys.CreateHDPath(viper.GetUint32(flagAccNumber), viper.GetUint32(flagAccIndex)).String(), cosmos.DefaultAlgo); err != nil {
				fmt.Println("keybase error")
				return err
			}

			// create rpc client
			rpc, err := cosmos.NewRPC(client, cdc, kb, cliCtx.ChainID, accName, accPass, viper.GetString(flagGasPrices))
			if err != nil {
				return err
			}

			// init event publisher
			ep := publisher.New(rpc)

			// orchestrator
			logger.Info("Starting orchestrator")
			orch := orchestrator.New(rpc, ep, logger)
			defer func() {
				logger.Info("Stopping orchestrator")
				orch.Stop()
			}()
			go func() {
				if err := orch.Start(); err != nil {
					logger.Error(err.Error())
					panic(err)
				}
			}()

			// init gRPC server.
			logger.Info("Starting gRPC server")
			server := grpc.New(rpc, ep, logger, viper.GetStringSlice(flagAuthorizedPubKeys))
			defer func() {
				logger.Info("Stopping gRPC server")
				server.Close()
			}()
			go func() {
				if err := server.Serve(viper.GetString(flagGrpcAddr)); err != nil {
					logger.Error(err.Error())
					panic(err)
				}
			}()

			<-xsignal.WaitForInterrupt()

			return nil
		},
	}
	cmd.Flags().String(flagGrpcAddr, ":50052", "The address for the gRPC server to expose")
	cmd.Flags().String(flagAuthorizedPubKeys, "", "The authorized pubkeys to communicate with the gRPC server")
	cmd.Flags().String(flagMnemonic, "", "The account's mnemonic that will be used to sign transactions")
	cmd.Flags().String(flagGasPrices, "", "The gas price to sign tx")
	cmd.Flags().String(flagAccNumber, "0", "The account number of the hd path to use to derive the mnemonic")
	cmd.Flags().String(flagAccIndex, "0", "The account index of the hd path to use to derive the mnemonic")
	return cmd
}

const (
	flagGrpcAddr          = "grpc-addr"
	flagAuthorizedPubKeys = "authorized-pubkeys"
	flagMnemonic          = "mnemonic"
	flagGasPrices         = "gas-prices"
	flagAccNumber         = "acc-number"
	flagAccIndex          = "acc-index"
)
