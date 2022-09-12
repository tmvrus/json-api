package entities

type BalanceResponse struct {
	Balance uint64 `json:"balance"`
}

type WithdrawAndDepositResponse struct {
	NewBalance uint64 `json:"newBalance"`
	TxID       string `json:"transactionId"`
}
