package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTransaction(t *testing.T) {

	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("AAAAAAAA")
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 2
	amount := int64(1)

	errs := make(chan error)
	results := make(chan TransferTransactionResult)

	transactions_dictionary := make(map[int]bool)

	for i := 0; i < n; i++ {
		transactionName := fmt.Sprintf("transaction %d", i+1)
		go func() {
			context_iteration := context.WithValue(context.Background(), transactionKey, transactionName)
			result, err := store.TransferTransaction(context_iteration, TransferTransactionParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Ammount:       amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err_iteration := <-errs
		require.NoError(t, err_iteration)

		result_iteration := <-results
		require.NotEmpty(t, result_iteration)

		//check transfers

		transfer := result_iteration.Transfer
		require.NotEmpty(t, transfer)

		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err_iteration = store.GetTransfer(context.Background(), transfer.ID)

		require.NoError(t, err_iteration)

		//check from entry

		fromEntry := result_iteration.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err_iteration = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err_iteration)

		//check to entry

		toEntry := result_iteration.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err_iteration = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err_iteration)

		// check FromAccounts

		fromAccount := result_iteration.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		// check ToAccounts

		toAccount := result_iteration.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println("transaction result: ", fromAccount.Balance, toAccount.Balance)

		diffAccount1 := account1.Balance - fromAccount.Balance
		diffAccount2 := toAccount.Balance - account2.Balance
		require.Equal(t, diffAccount1, diffAccount2)
		require.True(t, diffAccount1 > 0)
		require.True(t, diffAccount2 > 0)

		k := int(diffAccount1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, transactions_dictionary, k)
		transactions_dictionary[k] = true

	}

	// Update account balance
	updatedAccount1, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
