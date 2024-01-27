package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acct1 := createRandomAccount(t)
	acct2, err := testQueries.GetAccountForUpdate(context.Background(), acct1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acct2)

	require.Equal(t, acct1.ID, acct2.ID)
	require.Equal(t, acct1.Owner, acct2.Owner)
	require.Equal(t, acct1.Balance, acct2.Balance)
	require.Equal(t, acct1.Currency, acct2.Currency)
	require.WithinDuration(t, acct1.CreatedAt, acct2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acct1 := createRandomAccount(t)
	var updtParams = UpdateAccountParams{acct1.ID, 2000}
	acct2, err := testQueries.UpdateAccount(context.Background(), updtParams)
	require.NoError(t, err)
	require.NotEmpty(t, acct2)

	require.Equal(t, acct1.ID, acct2.ID)
	require.Equal(t, acct1.Owner, acct2.Owner)
	require.Equal(t, updtParams.Balance, acct2.Balance)
	require.Equal(t, acct1.Currency, acct2.Currency)
	require.WithinDuration(t, acct1.CreatedAt, acct2.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}
	listAccountParams := ListAccountParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccount(context.Background(), listAccountParams)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, _ := testQueries.GetAccountForUpdate(context.Background(), account.ID)
	require.Empty(t, account2)
}
