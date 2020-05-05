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

			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "orchestrator")
			client, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			// init rpc client
			logger.Info("starting rpc client")
			if err := client.Start(); err != nil {
				return err
			}
			defer func() {
				logger.Info("stopping rpc client")
				if err := client.Stop(); err != nil {
					logger.Error(err.Error())
				}
			}()

			kb := cosmos.NewInMemoryKeybase()
			if _, err := kb.CreateAccount("orchestrator", viper.GetString(flagMnemonic), "", "password", keys.CreateHDPath(0, 0).String(), cosmos.DefaultAlgo); err != nil {
				fmt.Println("keybase error")
				return err
			}

			// create rpc client
			rpc, err := cosmos.NewRPC(client, cdc, kb, cliCtx.ChainID, "orchestrator", "password", viper.GetString(flagGasPrices))
			if err != nil {
				return err
			}

			// init event publisher
			ep := publisher.New(rpc)

			// init gRPC server.
			logger.Info("starting grpc server")
			server := grpc.New(rpc, ep, viper.GetStringSlice(flagAuthorizedPubKeys))
			defer func() {
				logger.Info("stopping grpc server")
				server.Close()
			}()

			go func() {
				if err := server.Serve(viper.GetString(flagGrpcAddr)); err != nil {
					logger.Error(err.Error())
					panic(err)
				}
			}()

			// orchestrator
			logger.Info("starting orchestrator")
			orch := orchestrator.New(rpc, ep, viper.GetString(flagExecPrice))
			defer func() {
				logger.Info("stopping orchestrator")
				orch.Stop()
			}()
			go func() {
				if err := orch.Start(); err != nil {
					logger.Error(err.Error())
					panic(err)
				}
			}()
			go func() {
				for err := range orch.ErrC {
					logger.Error(err.Error())
				}
			}()

			<-xsignal.WaitForInterrupt()

			return nil
		},
	}
	cmd.Flags().String(flagGrpcAddr, ":50052", "The address for the gRPC server to expose")
	cmd.Flags().String(flagAuthorizedPubKeys, "", "The authorized pubkeys to communicate with the gRPC server")
	cmd.Flags().String(flagMnemonic, "", "The account's mnemonic that will be used to sign transactions")
	cmd.Flags().String(flagGasPrices, "1.0atto", "The gas price to sign tx")
	cmd.Flags().String(flagExecPrice, "10000atto", "The execution price to create execution")
	return cmd
}

const (
	flagGrpcAddr          = "grpc-addr"
	flagAuthorizedPubKeys = "authorized-pubkeys"
	flagMnemonic          = "mnemonic"
	flagGasPrices         = "gas-prices"
	flagExecPrice         = "exec-price"
)
