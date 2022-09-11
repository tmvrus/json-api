package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRollbackTransactionRequest_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		in    RollbackTransactionRequest
		error string
	}{
		{
			name: "happy path",
			in: RollbackTransactionRequest{
				CommonRequest: CommonRequest{
					CallerID:   1,
					PlayerName: "Player",
				},
				TxRef: "tx1",
			},
		},
		{
			name: "invalid tx ref",
			in: RollbackTransactionRequest{
				CommonRequest: CommonRequest{
					CallerID:   1,
					PlayerName: "Player",
				},
				TxRef: "",
			},
			error: "transactionRef not passed or invalid",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.in.Validate()
			if test.error == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.error)
			}
		})
	}
}

func TestWithdrawAndDepositRequest_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		in    WithdrawAndDepositRequest
		error string
	}{
		{
			name: "happy path",
			in: WithdrawAndDepositRequest{
				CommonRequest: CommonRequest{
					CallerID:   1,
					PlayerName: "Player",
				},
				Currency: "USD",
				TxRef:    "tx1",
				Withdraw: 10,
				Deposit:  20,
			},
		},
		{
			name: "invalid amount",
			in: WithdrawAndDepositRequest{
				CommonRequest: CommonRequest{
					CallerID:   1,
					PlayerName: "Player",
				},
				Currency: "USD",
				TxRef:    "tx1",
				Withdraw: 0,
				Deposit:  0,
			},
			error: "meaningless request: deposit and withdraws is zero",
		},
		{
			name: "invalid txRef",
			in: WithdrawAndDepositRequest{
				CommonRequest: CommonRequest{
					CallerID:   1,
					PlayerName: "Player",
				},
				Currency: "USD",
				TxRef:    "",
				Withdraw: 10,
				Deposit:  20,
			},
			error: "transactionRef not passed or invalid",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.in.Validate()
			if test.error == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.error)
			}
		})
	}
}

func TestBalanceRequest_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		in    BalanceRequest
		error string
	}{
		{
			name: "happy path",
			in: BalanceRequest{
				CommonRequest: CommonRequest{
					CallerID:   1,
					PlayerName: "Player",
				},
				Currency: "USD",
			},
		},
		{
			name: "Invalid player",
			in: BalanceRequest{
				CommonRequest: CommonRequest{
					CallerID:   1,
					PlayerName: "",
				},
				Currency: "USD",
			},
			error: "playerName not passed or invalid",
		},
		{
			name: "Invalid currency",
			in: BalanceRequest{
				CommonRequest: CommonRequest{
					CallerID:   1,
					PlayerName: "player2",
				},
				Currency: "?",
			},
			error: "ISO 4217 code",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.in.Validate()
			if test.error == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.error)
			}
		})
	}
}
