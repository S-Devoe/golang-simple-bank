package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// test create entry
func TestCreateEntry(t *testing.T) {

	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    100.0,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

}
func TestGetAllEntries(t *testing.T) {

	entries, err := testQueries.ListEntries(context.Background())
	require.NoError(t, err)
	require.Greater(t, len(entries), 0)

}

// test get entry
func TestGetEntry(t *testing.T) {

	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    100.0,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)

	entryGot, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.Equal(t, entryGot.ID, entry.ID)
	require.Equal(t, entryGot.AccountID, entry.AccountID)
	require.Equal(t, entryGot.Amount, entry.Amount)

}
