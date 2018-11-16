package cli

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/cosmos/cosmos-sdk/x/bank/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
	flagNums   = "nums"
	flagPack   = "pack"
	flagStep   = "step"
)

type sendPack struct {
	From string `json:"from"`
	To   string `json:"to"`
	//amount string `json:"amount"`
}

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create and sign a send tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			amount := viper.GetString(flagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// ensure account has enough coins
			if !account.GetCoins().IsGTE(coins) {
				return errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := client.BuildMsg(from, to, coins)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")

	return cmd
}

// BatchSendTxCmd will create send txs and sign them with the given key.
func BatchSendTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-send",
		Short: "Create and sign send txs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)
			_ = toStr

			// parse coins trying to be sent
			amount := viper.GetString(flagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			// parse nums of txs to be sent
			numsStr := viper.GetString(flagNums)
			nums, err := strconv.ParseInt(numsStr, 10, 64)
			if err != nil || nums == 0 {
				return err
			}

			// parse nums of txs to be sent
			stepStr := viper.GetString(flagStep)
			//step, err := strconv.Atoi(stepStr)
			step, err := strconv.ParseInt(stepStr, 10, 64)
			if err != nil || step == 0 {
				return err
			}

			passphrase, err := keys.GetPassphrase(cliCtx.FromAddressName)
			if err != nil {
				return err
			}

			var sendArray []sendPack
			packStr := viper.GetString(flagPack)
			err = json.Unmarshal([]byte(packStr), &sendArray)
			if err != nil || len(sendArray) == 0 {
				return err
			}

			var wg sync.WaitGroup
			wg.Add(len(sendArray))

			for _, key := range sendArray {
				go sendAccountTxs(cdc, key.From, key.To, coins, nums, step, passphrase, &wg)
			}

			wg.Wait()
			return nil
		},
	}

	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")
	cmd.Flags().String(flagNums, "", "nums of txs to send")
	cmd.Flags().String(flagPack, "", "packet of msgs to send")
	cmd.Flags().String(flagStep, "", "numbers of txs in one step")

	return cmd
}

func sendAccountTxs(cdc *wire.Codec, fromStr, toStr string, coins sdk.Coins, nums, step int64, passphrase string, wg *sync.WaitGroup) error {
	txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
	cliCtx := context.NewCLIContext().
		WithCodec(cdc).
		WithLogger(os.Stdout).
		WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

	// get from account later
	txCtx.AccountNumber = 0
	txCtx.Sequence = 0
	cliCtx.FromAddressName = fromStr

	if err := cliCtx.EnsureAccountExists(); err != nil {
		return err
	}

	from, err := cliCtx.GetFromAddress()
	if err != nil {
		return err
	}

	to, err := sdk.AccAddressFromBech32(toStr)
	if err != nil {
		return err
	}

	account, err := cliCtx.GetAccount(from)
	if err != nil {
		return err
	}

	// ensure account has enough coins
	if !account.GetCoins().IsGTE(coins) {
		return errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := client.BuildMsg(from, to, coins)

	err = utils.BatchSendTx(txCtx, cliCtx, []sdk.Msg{msg}, nums, step, passphrase)

	wg.Done()

	return err
}
