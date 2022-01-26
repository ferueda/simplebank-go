package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))
}

func TestGetTransfer(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)

	transfer := createRandomTransfer(t, fromAcc, toAcc)

	queriedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, queriedTransfer)
	require.Equal(t, transfer.ID, queriedTransfer.ID)
	require.Equal(t, transfer.CreatedAt, queriedTransfer.CreatedAt)
	require.Equal(t, transfer.FromAccountID, queriedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, queriedTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, queriedTransfer.Amount)
}

func TestListTransfer(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, fromAcc, toAcc)
		createRandomTransfer(t, toAcc, fromAcc)
	}

	arg := ListTransfersParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   fromAcc.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == fromAcc.ID || transfer.ToAccountID == fromAcc.ID)
	}
}

func createRandomTransfer(t *testing.T, from, to Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        randomInt(100, 1_000_000),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	return transfer
}
