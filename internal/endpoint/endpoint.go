package endpoint

import (
	"context"

	"github.com/tmvrus/json-api/internal/entities"
	"github.com/tmvrus/json-api/internal/storage"
)

type Endpoint struct {
	balance storage.BalanceStorage
}

func NewEndpoint(s storage.BalanceStorage) *Endpoint {
	return &Endpoint{balance: s}
}

func (e *Endpoint) getBalance(ctx context.Context, r entities.BalanceRequest) {

}

func (e *Endpoint) rollbackTransaction(ctx context.Context, r entities.BalanceRequest) {

}

func (e *Endpoint) withdrawAndDeposit() {

}
