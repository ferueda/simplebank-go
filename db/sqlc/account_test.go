package db

import (
	"context"
	"database/sql"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createdAcc := createRandomAccount(t)
	queriedAcc, err := testQueries.GetAccount(context.Background(), createdAcc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, queriedAcc)
	require.Equal(t, queriedAcc.ID, createdAcc.ID)
	require.Equal(t, queriedAcc.Owner, createdAcc.Owner)
	require.Equal(t, queriedAcc.Balance, createdAcc.Balance)
	require.Equal(t, queriedAcc.CreatedAt, createdAcc.CreatedAt)
	require.Equal(t, queriedAcc.Currency, createdAcc.Currency)
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
	require.Equal(t, updatedAcc.ID, createdAcc.ID)
	require.Equal(t, updatedAcc.Owner, createdAcc.Owner)
	require.Equal(t, updatedAcc.Balance, arg.Balance)
	require.Equal(t, updatedAcc.CreatedAt, createdAcc.CreatedAt)
	require.Equal(t, updatedAcc.Currency, createdAcc.Currency)
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

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
