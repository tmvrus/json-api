package storage

import (
	"context"

	"github.com/tmvrus/json-api/internal/entities"
)

type BalanceStorage interface {
	GetBalance(ctx context.Context, r entities.BalanceRequest) (entities.BalanceResponse, error)
	RollbackTransaction(_ context.Context, txRef string) error
	WithdrawAndDeposit(_ context.Context, r entities.WithdrawAndDepositRequest) (entities.WithdrawAndDepositResponse, error)
}
