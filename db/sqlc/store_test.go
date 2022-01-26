package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	s := NewStore(testDB)

	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)

	n := 5
	transferAmount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := s.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Amount:        transferAmount,
			})

			errors <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		require.NotEmpty(t, result.Transfer)
		require.Equal(t, fromAcc.ID, result.Transfer.FromAccountID)
		require.Equal(t, toAcc.ID, result.Transfer.ToAccountID)
		require.Equal(t, transferAmount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = s.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, fromAcc.ID, result.FromEntry.AccountID)
		require.Equal(t, -transferAmount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)

		_, err = s.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, toAcc.ID, result.ToEntry.AccountID)
		require.Equal(t, transferAmount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)

		_, err = s.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, fromAcc.ID, result.FromAccount.ID)
		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, toAcc.ID, result.ToAccount.ID)

		diffFromAcc := fromAcc.Balance - result.FromAccount.Balance
		diffToAcc := result.ToAccount.Balance - toAcc.Balance
		require.Equal(t, diffFromAcc, diffToAcc)
		require.True(t, diffFromAcc > 0)
		require.True(t, diffFromAcc%transferAmount == 0)
	}

	updatedFromAcc, err := testQueries.GetAccount(context.Background(), fromAcc.ID)
	require.NoError(t, err)

	updatedToAcc, err := testQueries.GetAccount(context.Background(), toAcc.ID)
	require.NoError(t, err)

	require.Equal(t, fromAcc.Balance-int64(n)*transferAmount, updatedFromAcc.Balance)
	require.Equal(t, toAcc.Balance+int64(n)*transferAmount, updatedToAcc.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	s := NewStore(testDB)

	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)

	n := 10
	transferAmount := int64(10)

	errors := make(chan error)

	for i := 0; i < n; i++ {
		fromAccId := fromAcc.ID
		toAccId := toAcc.ID

		if i%2 == 1 {
			fromAccId = toAcc.ID
			toAccId = fromAcc.ID
		}

		go func() {
			_, err := s.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccId,
				ToAccountID:   toAccId,
				Amount:        transferAmount,
			})

			errors <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	updatedFromAcc, err := testQueries.GetAccount(context.Background(), fromAcc.ID)
	require.NoError(t, err)

	updatedToAcc, err := testQueries.GetAccount(context.Background(), toAcc.ID)
	require.NoError(t, err)

	require.Equal(t, fromAcc.Balance, updatedFromAcc.Balance)
	require.Equal(t, toAcc.Balance, updatedToAcc.Balance)
}
