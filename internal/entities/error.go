package entities

import "fmt"

var (
	ErrBalanceNotEnough = fmt.Errorf("balance is not enough")
	ErrTxWasRolledBack  = fmt.Errorf("transaction was rolled back")
)
