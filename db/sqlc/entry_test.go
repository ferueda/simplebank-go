package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t, createRandomAccount(t))
}

func TestGetEntry(t *testing.T) {
	createdEntry := createRandomEntry(t, createRandomAccount(t))
	queriedEntry, err := testQueries.GetEntry(context.Background(), createdEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, queriedEntry)
	require.Equal(t, createdEntry.ID, queriedEntry.ID)
	require.Equal(t, createdEntry.AccountID, queriedEntry.AccountID)
	require.Equal(t, createdEntry.Amount, queriedEntry.Amount)
	require.Equal(t, createdEntry.CreatedAt, queriedEntry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	entry := createRandomEntry(t, createRandomAccount(t))

	arg := ListEntriesParams{
		AccountID: entry.AccountID,
		Limit:     5,
		Offset:    0,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, entry.AccountID, entries[0].AccountID)
	require.Equal(t, entry.Amount, entries[0].Amount)
	require.Equal(t, entry.CreatedAt, entries[0].CreatedAt)
	require.Equal(t, entry.ID, entries[0].ID)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t, createRandomAccount(t))
	err := testQueries.DeleteEntry(context.Background(), entry.AccountID)
	require.NoError(t, err)

	queriedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, queriedEntry)
}

func createRandomEntry(t *testing.T, acc Account) Entry {
	arg := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    randomInt(1, 1_000_000),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}
