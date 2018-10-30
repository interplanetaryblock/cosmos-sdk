package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/cosmos/cosmos-sdk/examples/multicoin/app"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "multicoind",
		Short:             "Multicoin Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, server.DefaultAppInit,
		server.ConstructAppCreator(newApp, "multicoin"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "multicoin"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.multicoind")
	executor := cli.PrepareBaseCmd(rootCmd, "MC", rootDir)

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	return app.NewMulticoinApp(logger, db, baseapp.SetPruning(viper.GetString("pruning")))
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewMulticoinApp(logger, db)
	return bapp.ExportAppStateAndValidators()
}
