package utils

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
)

// SendTx implements a auxiliary handler that facilitates sending a series of
// messages in a signed transaction given a TxContext and a QueryContext. It
// ensures that the account exists, has a proper number and sequence set. In
// addition, it builds and signs a transaction with the supplied messages.
// Finally, it broadcasts the signed transaction to a node.
func SendTx(txCtx authctx.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) error {
	if err := cliCtx.EnsureAccountExists(); err != nil {
		return err
	}

	from, err := cliCtx.GetFromAddress()
	if err != nil {
		return err
	}

	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if txCtx.AccountNumber == 0 {
		accNum, err := cliCtx.GetAccountNumber(from)
		if err != nil {
			return err
		}

		txCtx = txCtx.WithAccountNumber(accNum)
	}

	// TODO: (ref #1903) Allow for user supplied account sequence without
	// automatically doing a manual lookup.
	if txCtx.Sequence == 0 {
		accSeq, err := cliCtx.GetAccountSequence(from)
		if err != nil {
			return err
		}

		txCtx = txCtx.WithSequence(accSeq)
	}

	passphrase, err := keys.GetPassphrase(cliCtx.FromAddressName)
	if err != nil {
		return err
	}

	// build and sign the transaction
	txBytes, err := txCtx.BuildAndSign(cliCtx.FromAddressName, passphrase, msgs)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	return cliCtx.EnsureBroadcastTx(txBytes)
}

// BatchSendTx implements a auxiliary handler that facilitates sending a series of
// messages in a signed transaction given a TxContext and a QueryContext. It
// ensures that the account exists, has a proper number and sequence set. In
// addition, it builds and signs a transaction with the supplied messages.
// Finally, it broadcasts the signed transaction to a node.
func BatchSendTx(txCtx authctx.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg, nums, step int64, passphrase string) error {
	if err := cliCtx.EnsureAccountExists(); err != nil {
		return err
	}

	from, err := cliCtx.GetFromAddress()
	if err != nil {
		return err
	}

	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	accNum, err := cliCtx.GetAccountNumber(from)
	if err != nil {
		return err
	}

	txCtx = txCtx.WithAccountNumber(accNum)

	// TODO: (ref #1903) Allow for user supplied account sequence without
	// automatically doing a manual lookup.
	accSeq, err := cliCtx.GetAccountSequence(from)
	if err != nil {
		return err
	}

	txCtx = txCtx.WithSequence(accSeq)

	// using the passphrase args
	//passphrase, err := keys.GetPassphrase(cliCtx.FromAddressName)
	//if err != nil {
	//	return err
	//}

	txMap := make(map[int64][]byte)
	total := accSeq + nums

	// batch build and sign txs
	for i := accSeq; i < total; i++ {
		// build and sign the transaction
		txBytes, err := txCtx.BuildAndSign(cliCtx.FromAddressName, passphrase, msgs)
		if err != nil {
			return err
		}

		txMap[i] = txBytes
		fmt.Printf("account number: %d, sign %d msg, size: %d, progress:%d%%\n",
			txCtx.AccountNumber, i, len(txBytes), (i-accSeq)*100/nums)

		txCtx.Sequence++
	}

	// count nums of txs has been sent
	count := int64(0)
	curSeq := accSeq

	// batch broadcast txs
	for i := accSeq; i < total; i++ {
		if count > 0 && count%step == 0 {
			// get current seq
			curSeq, err = cliCtx.GetAccountSequence(from)
			if err != nil {
				return err
			}

			if curSeq < accSeq || curSeq >= total {
				return fmt.Errorf("wrong seq, should be: %d, cur: %d", i, curSeq)
			}

			//if curSeq != accSeq && curSeq != i {
			//fmt.Printf("warnning: seq should be: %d, cur: %d", i, curSeq)
			//}
			time.Sleep(time.Duration(10) * time.Millisecond)
		}

		// broadcast to a Tendermint node
		err = cliCtx.EnsureBroadcastTx(txMap[i])
		if err != nil {
			fmt.Println(err)
		}

		count++
		fmt.Printf("account number: %d, broadcast %d , progress:%d%%\n",
			txCtx.AccountNumber, i, count*100/nums)

		// send txs no more than nums
		if count == nums {
			return nil
		}
	}

	return nil
}
