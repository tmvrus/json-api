package entities

type BalanceRequest struct {
	CallerID   int64    `json:"callerId"`
	PlayerName string   `json:"playerName"`
	Currency   Currency `json:"currency"`
}

type WithdrawAndDepositRequest struct {
	CallerID   int64    `json:"callerId"`
	PlayerName string   `json:"playerName"`
	Currency   Currency `json:"currency"`
	TxRef      string   `json:"transactionRef"`
	Withdraw   uint64   `json:"withdraw"`
	Deposit    uint64   `json:"deposit"`
}
