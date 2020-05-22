package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/gorilla/mux"
	"github.com/mesg-foundation/engine/app"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"
)

// ServeCommand creates and starts the LCD server
// adapted version of function from https://github.com/cosmos/cosmos-sdk/blob/v0.38.3/client/lcd/root.go#L74-L100
func ServeCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// new rest server
			r := mux.NewRouter()
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			logger := log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")

			// register routes
			client.RegisterRoutes(cliCtx, r)
			authrest.RegisterTxRoutes(cliCtx, r)
			app.ModuleBasics.RegisterRESTRoutes(cliCtx, r)
			cosmos.RegisterSimulateRoute(cliCtx, r)

			// start
			var listener net.Listener
			server.TrapSignal(func() {
				err := listener.Close()
				logger.Error("error closing listener", "err", err)
			})

			cfg := rpcserver.DefaultConfig()
			cfg.MaxOpenConnections = viper.GetInt(flags.FlagMaxOpenConnections)
			cfg.ReadTimeout = time.Duration(uint(viper.GetInt(flags.FlagRPCReadTimeout))) * time.Second
			cfg.WriteTimeout = time.Duration(uint(viper.GetInt(flags.FlagRPCWriteTimeout))) * time.Second

			listener, err = rpcserver.Listen(viper.GetString(flags.FlagListenAddr), cfg)
			if err != nil {
				return
			}
			logger.Info(
				fmt.Sprintf(
					"Starting application REST service (chain-id: %q)...",
					viper.GetString(flags.FlagChainID),
				),
			)

			return rpcserver.StartHTTPServer(listener, r, logger, cfg)
		},
	}

	return flags.RegisterRestServerFlags(cmd)
}
