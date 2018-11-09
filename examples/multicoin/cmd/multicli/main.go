package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	govcmd "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	ibccmd "github.com/cosmos/cosmos-sdk/x/ibc/client/cli"
	namingcmd "github.com/cosmos/cosmos-sdk/x/naming/client/cli"
	stakecmd "github.com/cosmos/cosmos-sdk/x/stake/client/cli"

	"github.com/cosmos/cosmos-sdk/examples/multicoin/app"
	"github.com/cosmos/cosmos-sdk/examples/multicoin/types"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "multicli",
		Short: "Multicoin light-client",
	}
)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()

	// TODO: Setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc.

	// add standard rpc, and tx commands
	rpc.AddCommands(rootCmd)
	rootCmd.AddCommand(client.LineBreak)
	tx.AddCommands(rootCmd, cdc)
	rootCmd.AddCommand(client.LineBreak)

	// add query/post commands (custom to binary)
	rootCmd.AddCommand(
		client.GetCommands(
			stakecmd.GetCmdQueryValidator("stake", cdc),
			stakecmd.GetCmdQueryValidators("stake", cdc),
			stakecmd.GetCmdQueryDelegation("stake", cdc),
			stakecmd.GetCmdQueryDelegations("stake", cdc),
			authcmd.GetAccountCmd("acc", cdc, types.GetAccountDecoder(cdc)),
		)...)

	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
			ibccmd.IBCTransferCmd(cdc),
			ibccmd.IBCRelayCmd(cdc),
			stakecmd.GetCmdCreateValidator(cdc),
			stakecmd.GetCmdEditValidator(cdc),
			stakecmd.GetCmdDelegate(cdc),
			stakecmd.GetCmdUnbond("stake", cdc),
		)...)

	// add gov commands
	rootCmd.AddCommand(
		client.GetCommands(
			govcmd.GetCmdQueryProposal("gov", cdc),
			govcmd.GetCmdQueryProposals("gov", cdc),
			govcmd.GetCmdQueryVote("gov", cdc),
			govcmd.GetCmdQueryVotes("gov", cdc),
		)...)
	rootCmd.AddCommand(client.PostCommands(
		govcmd.GetCmdSubmitProposal(cdc),
		govcmd.GetCmdDeposit(cdc),
		govcmd.GetCmdVote(cdc),
	)...)

	// add naming commands
	rootCmd.AddCommand(
		client.GetCommands(
			namingcmd.GetCmdResolveName("naming", cdc),
			namingcmd.GetCmdWhois("naming", cdc),
		)...)
	rootCmd.AddCommand(
		client.PostCommands(
			namingcmd.GetCmdBuyName(cdc),
			namingcmd.GetCmdSetName(cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "MC", os.ExpandEnv("$HOME/.multicli"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}
