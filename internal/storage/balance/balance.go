package balance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tmvrus/json-api/internal/entities"
)

type Storage struct {
	balanceLock sync.RWMutex
	balanceData map[string]uint64

	txLock sync.RWMutex
	txLog  map[string]txInfo
}

func NewStorage() *Storage {
	return &Storage{
		balanceLock: sync.RWMutex{},
		txLock:      sync.RWMutex{},
		balanceData: map[string]uint64{},
		txLog:       map[string]txInfo{},
	}
}

func (s *Storage) GetBalance(_ context.Context, r entities.BalanceRequest) (entities.BalanceResponse, error) {
	key := fmt.Sprintf("%s-%s", r.PlayerName, r.Currency)

	s.balanceLock.RLock()
	v := s.balanceData[key]
	s.balanceLock.RUnlock()

	return entities.BalanceResponse{
		Balance: v,
	}, nil
}

func (s *Storage) RollbackTransaction(_ context.Context, txRef string) error {
	s.txLock.Lock()
	defer s.txLock.Unlock()

	v, ok := s.txLog[txRef]
	if !ok {
		s.txLog[txRef] = txInfo{
			Timestamp: time.Now().UTC(),
			Status:    txStatusRollback,
		}
		return nil
	}

	if v.Status == txStatusRollback {
		return nil
	}

	s.balanceLock.Lock()
	defer s.balanceLock.Unlock()

	b, ok := s.balanceData[v.BalanceRef]
	if !ok {
		return fmt.Errorf("failed to rollaback trasaction %s: balance record was not found", txRef)
	}

	b = b + v.Add
	b = b - v.Sub
	s.balanceData[v.BalanceRef] = b

	v.Status = txStatusRollback
	s.txLog[txRef] = v

	return nil
}

func (s *Storage) WithdrawAndDeposit(_ context.Context, r entities.WithdrawAndDepositRequest) (entities.WithdrawAndDepositResponse, error) {
	s.txLock.Lock()
	defer s.txLock.Unlock()
	s.balanceLock.Lock()
	defer s.balanceLock.Unlock()

	var res entities.WithdrawAndDepositResponse

	tx, ok := s.txLog[r.TxRef]
	if ok {
		switch tx.Status {
		case txStatusRollback:
			return res, fmt.Errorf("failed to process WithdrawAndDepositRequest for tx %s: %w", r.TxRef, entities.ErrTxWasRolledBack)
		case txStatusCommit:
			return res, nil
		default:
			return res, fmt.Errorf("failed to process WithdrawAndDepositRequest: balanceData corruption detected, invalid tx status")
		}
	}

	balanceRef := fmt.Sprintf("%s-%s", r.PlayerName, r.Currency)
	originBalance := s.balanceData[balanceRef]
	if r.Withdraw > originBalance {
		return res, fmt.Errorf("failed to process WithdrawAndDepositRequest: %w: (%d) is less than withdraw (%d)", entities.ErrBalanceNotEnough, originBalance, r.Withdraw)
	}

	tx = txInfo{
		Status:     txStatusCommit,
		Timestamp:  time.Now().UTC(),
		BalanceRef: balanceRef,
	}

	newBalance := originBalance - r.Withdraw + r.Deposit
	if newBalance > originBalance {
		tx.Sub = newBalance - originBalance
	} else {
		tx.Add = originBalance - newBalance
	}

	s.balanceData[balanceRef] = newBalance
	s.txLog[r.TxRef] = tx

	res.NewBalance = newBalance
	res.TxID = r.TxRef

	return res, nil
}
