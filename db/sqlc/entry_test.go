package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T, account *Account, amount int64) Entry {
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount: amount,
	}
	entry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.Amount, entry.Amount)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.NotZero(t, entry.ID)
	require.False(t, entry.CreatedAt.IsZero())

	return entry
}

func TestCreateEntry(t *testing.T) {
	defer cleanup()

	account := createRandomAccount(t)
	createTestEntry(t, &account, 100)
}

func TestGetEntry(t *testing.T) {
	defer cleanup()

	account := createRandomAccount(t)
	entry1 := createTestEntry(t, &account, 100)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
}

func TestListEntries(t *testing.T) {
	defer cleanup()

	account := createRandomAccount(t)
	for i := 0; i < 3; i++ {
		createTestEntry(t, &account, 100)
	}

	args := ListEntriesParams {
		AccountID: account.ID,
		Limit: 3,
		Offset: 0,
	}
	entries2, err := testQueries.ListEntries(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t,entries2)

	require.Len(t, entries2, 3)
	for _, entry := range entries2 {
		require.NotEmpty(t, entry)
	}
}


