package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/stretchr/testify/require"
)

// instead of typing the code for creating account in each test all the time, i can just call this
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.GenerateRandomName(),
		Balance:  util.RandomMoney(),
		Currency: util.GenerateRandomCurrency(),
	}

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
	arg := CreateAccountParams{
		Owner:    util.GenerateRandomName(),
		Balance:  util.RandomMoney(),
		Currency: util.GenerateRandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}

func TestGetAccount(t *testing.T) {
	createArg := CreateAccountParams{
		Owner:    util.GenerateRandomName(),
		Balance:  util.RandomMoney(),
		Currency: util.GenerateRandomCurrency(),
	}

	account1, err := testQueries.CreateAccount(context.Background(), createArg)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 0)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account1.ID, account.ID)
	require.Equal(t, account1.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, account1.Currency, account.Currency)

}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	acct, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.NotEqual(t, account.ID, acct.ID)
	require.Empty(t, acct)

}

func TestListAccounts(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
