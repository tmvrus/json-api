package entities

import "fmt"

type CommonRequest struct {
	CallerID   int64  `json:"callerId"`
	PlayerName string `json:"playerName"`
}

type BalanceRequest struct {
	CommonRequest
	Currency Currency `json:"currency"`
}

type WithdrawAndDepositRequest struct {
	CommonRequest
	Currency Currency `json:"currency"`
	TxRef    string   `json:"transactionRef"`
	Withdraw uint64   `json:"withdraw"`
	Deposit  uint64   `json:"deposit"`
}

type RollbackTransactionRequest struct {
	CommonRequest
	TxRef string `json:"transactionRef"`
}

func (r *CommonRequest) Validate() error {
	if r.CallerID <= 0 {
		return fmt.Errorf("callerId not passed or invalid")
	}
	if r.PlayerName == "" {
		return fmt.Errorf("playerName not passed or invalid")
	}
	return nil
}

func (r *BalanceRequest) Validate() error {
	if err := r.Currency.Valid(); err != nil {
		return err
	}
	return r.CommonRequest.Validate()
}

func (r *WithdrawAndDepositRequest) Validate() error {
	if err := r.Currency.Valid(); err != nil {
		return err
	}
	if err := r.CommonRequest.Validate(); err != nil {
		return err
	}
	if r.TxRef == "" {
		return fmt.Errorf("transactionRef not passed or invalid")
	}
	if r.Withdraw == 0 && r.Deposit == 0 {
		return fmt.Errorf("meaningless request: deposit and withdraws is zero")
	}
	return nil
}

func (r *RollbackTransactionRequest) Validate() error {
	if err := r.CommonRequest.Validate(); err != nil {
		return err
	}
	if r.TxRef == "" {
		return fmt.Errorf("transactionRef not passed or invalid")
	}
	return nil
}
