package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/cosmos/cosmos-sdk/x/naming"
)

const (
	flagName   = "name"
	flagValue  = "value"
	flagAmount = "amount"
)

// GetCmdBuyName is the CLI command for sending a BuyName transaction
func GetCmdBuyName(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy-name [name] [amount]",
		Short: "bid for existing name or claim new name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			name := args[0]

			amount := args[1]
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			account, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := naming.NewMsgBuyName(name, coins, account)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Build and sign the transaction, then broadcast to Tendermint
			// NameID must be returned, and it is a part of response.
			cliCtx.PrintResponse = true
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdSetName is the CLI command for sending a SetName transaction
func GetCmdSetName(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-name [name] [value]",
		Short: "set the value associated with a name that you own",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			name := args[0]
			value := args[1]

			account, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := naming.NewMsgSetName(name, value, account)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
