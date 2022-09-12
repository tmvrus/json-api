package balance

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tmvrus/json-api/internal/entities"
)

func TestStorage_BalanceManagement(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s := NewStorage()

	_, err := s.WithdrawAndDeposit(ctx, entities.WithdrawAndDepositRequest{
		CommonRequest: entities.CommonRequest{
			PlayerName: "player1",
		},
		Currency: "USD",
		TxRef:    "tx0",
		Withdraw: 10,
		Deposit:  50,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "balance is not enough")

	s.balanceData["player1-USD"] = 10

	newBalance, err := s.WithdrawAndDeposit(ctx, entities.WithdrawAndDepositRequest{
		CommonRequest: entities.CommonRequest{
			PlayerName: "player1",
		},
		Currency: "USD",
		TxRef:    "tx1",
		Withdraw: 10,
		Deposit:  50,
	})
	require.NoError(t, err)
	require.Equal(t, newBalance, entities.WithdrawAndDepositResponse{
		NewBalance: 50,
		TxID:       "tx1",
	})

	balance, err := s.GetBalance(ctx, entities.BalanceRequest{
		CommonRequest: entities.CommonRequest{
			PlayerName: "player1",
		},
		Currency: "USD",
	})
	require.NoError(t, err)
	require.Equal(t, balance.Balance, uint64(50))

	_, err = s.WithdrawAndDeposit(ctx, entities.WithdrawAndDepositRequest{
		CommonRequest: entities.CommonRequest{
			PlayerName: "player1",
		},
		Currency: "USD",
		TxRef:    "tx1", // same txRef, nothing should change
		Withdraw: 10,
		Deposit:  50,
	})
	require.NoError(t, err)

	balance, err = s.GetBalance(ctx, entities.BalanceRequest{
		CommonRequest: entities.CommonRequest{
			PlayerName: "player1",
		},
		Currency: "USD",
	})
	require.NoError(t, err)
	require.Equal(t, balance.Balance, uint64(50))

	_, err = s.WithdrawAndDeposit(ctx, entities.WithdrawAndDepositRequest{
		CommonRequest: entities.CommonRequest{
			PlayerName: "player1",
		},
		Currency: "USD",
		TxRef:    "tx2",
		Withdraw: 20,
		Deposit:  30,
	})
	require.NoError(t, err)

	balance, err = s.GetBalance(ctx, entities.BalanceRequest{
		CommonRequest: entities.CommonRequest{
			PlayerName: "player1",
		},
		Currency: "USD",
	})
	require.NoError(t, err)
	require.Equal(t, balance.Balance, uint64(60))

	err = s.RollbackTransaction(ctx, "tx2")
	require.NoError(t, err)

	balance, err = s.GetBalance(ctx, entities.BalanceRequest{
		CommonRequest: entities.CommonRequest{
			PlayerName: "player1",
		},
		Currency: "USD",
	})
	require.NoError(t, err)
	require.Equal(t, balance.Balance, uint64(50))

	err = s.RollbackTransaction(ctx, "tx2")
	require.NoError(t, err)

	err = s.RollbackTransaction(ctx, "tx3")
	require.NoError(t, err)

	_, err = s.WithdrawAndDeposit(ctx, entities.WithdrawAndDepositRequest{
		CommonRequest: entities.CommonRequest{
			PlayerName: "player1",
		},
		Currency: "USD",
		TxRef:    "tx3",
		Withdraw: 10,
		Deposit:  5,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction was rolled back")
}
