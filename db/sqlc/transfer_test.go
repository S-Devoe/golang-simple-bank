package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// test create transfer
func TestCreateTransfer(t *testing.T) {
	// write the test for this  function
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	transfer, err := testStore.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, transfer.FromAccountID, account1.ID)
	require.Equal(t, transfer.ToAccountID, account2.ID)

}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	_, err := testStore.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)

	transfers, err := testStore.ListTransfers(context.Background())
	require.NoError(t, err)
	require.Equal(t, transfers[0].FromAccountID, account1.ID)
	require.Greater(t, len(transfers), 0)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	transfer, err := testStore.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)

	transferGot, err := testStore.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.Equal(t, transferGot.ID, transfer.ID)
	require.Equal(t, transferGot.FromAccountID, transfer.FromAccountID)
	require.Equal(t, transferGot.ToAccountID, transfer.ToAccountID)
}
