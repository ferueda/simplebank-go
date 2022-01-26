package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createdAcc := createRandomAccount(t)
	queriedAcc, err := testQueries.GetAccount(context.Background(), createdAcc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, queriedAcc)
	require.Equal(t, createdAcc.ID, queriedAcc.ID)
	require.Equal(t, createdAcc.Owner, queriedAcc.Owner)
	require.Equal(t, createdAcc.Balance, queriedAcc.Balance)
	require.Equal(t, createdAcc.CreatedAt, queriedAcc.CreatedAt)
	require.Equal(t, createdAcc.Currency, queriedAcc.Currency)
}

func TestUpdateAccount(t *testing.T) {
	createdAcc := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      createdAcc.ID,
		Balance: randomInt(0, 1_000_000),
	}

	updatedAcc, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)
	require.Equal(t, createdAcc.ID, updatedAcc.ID)
	require.Equal(t, createdAcc.Owner, updatedAcc.Owner)
	require.Equal(t, arg.Balance, updatedAcc.Balance)
	require.Equal(t, createdAcc.CreatedAt, updatedAcc.CreatedAt)
	require.Equal(t, createdAcc.Currency, updatedAcc.Currency)
}

func TestAddAccountBalance(t *testing.T) {
	createdAcc := createRandomAccount(t)

	arg := AddAccountBalanceParams{
		ID:     createdAcc.ID,
		Amount: randomInt(0, 1_000_000),
	}

	updatedAcc, err := testQueries.AddAccountBalance(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)
	require.Equal(t, createdAcc.ID, updatedAcc.ID)
	require.Equal(t, createdAcc.Owner, updatedAcc.Owner)
	require.Equal(t, createdAcc.Balance+arg.Amount, updatedAcc.Balance)
	require.Equal(t, createdAcc.CreatedAt, updatedAcc.CreatedAt)
	require.Equal(t, createdAcc.Currency, updatedAcc.Currency)
}

func TestDeleteAccount(t *testing.T) {
	createdAcc := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), createdAcc.ID)

	require.NoError(t, err)

	queriedAcc, err := testQueries.GetAccount(context.Background(), createdAcc.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, queriedAcc)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    randomString(6),
		Balance:  randomInt(0, 1_000_000),
		Currency: "CAD",
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
