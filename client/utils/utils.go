package utils

import (
	"fmt"

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
func BatchSendTx(txCtx authctx.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg, nums int) error {
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

	txMap := make(map[int][]byte)

	// batch build and sign txs
	for i := 0; i < nums; i++ {
		// build and sign the transaction
		txBytes, err := txCtx.BuildAndSign(cliCtx.FromAddressName, passphrase, msgs)
		if err != nil {
			return err
		}

		txMap[i] = txBytes
		fmt.Printf("account number: %d, sign %d msg, size: %d\n", txCtx.AccountNumber, i, len(txBytes))

		txCtx.Sequence++
	}

	// batch broadcast txs
	for i := 0; i < nums; i++ {
		// broadcast to a Tendermint node
		err = cliCtx.EnsureBroadcastTx(txMap[i])
		if err != nil {
			return err
		}

		fmt.Printf("account number: %d, broadcast %d tx:\n", txCtx.AccountNumber, i)
	}

	return nil
}
