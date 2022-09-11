package balance

import (
	"time"
)

type txStatus uint

const (
	txStatusInvalid txStatus = iota
	txStatusCommit
	txStatusRollback
)

type txInfo struct {
	Timestamp  time.Time
	Status     txStatus
	BalanceRef string
	Sub        uint64
	Add        uint64
}
