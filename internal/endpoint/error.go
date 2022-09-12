package endpoint

const (
	codeInvalidJSON         int64 = -32600
	codeMethodNotFound            = -32601
	codeInvalidMethodParams       = -32602
	codeInternalError             = -32603
	codeBalanceNotEnough          = -32001
	codeTxWasRolledBack           = -32002
)

type apiError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
