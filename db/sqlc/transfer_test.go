package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/MasonPhan2110/SimpleBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	return transfer
}

func TestCreateTrasnfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestGetListTransfer(t *testing.T) {
	for i := 0; i < 15; i++ {
		createRandomTransfer(t)
	}
	arg := ListTransferParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Limit:         10,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfer(context.Background(), arg)
	log.Println(transfers)
	require.NoError(t, err)
	require.Len(t, transfers, 0)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
