package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	utils "github.com/quiless/BANK_GO/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {

	arg := CreateAccountParams{
		Owner:    utils.GenerateString(6),
		Balance:  utils.GenerateInt(0, 1000),
		Currency: utils.GetCurrencyType()}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account

}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccoun(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Millisecond)

}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.GenerateInt(500, 1000),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Millisecond)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T) {

	createRandomAccount(t)

	arg := ListAccountsParams{
		Limit:  1,
		Offset: 1,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 1)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
