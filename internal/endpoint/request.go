package endpoint

import (
	"encoding/json"
	"fmt"
)

type baseRequest struct {
	Version   string          `json:"jsonrpc"`
	Method    string          `json:"method"`
	RequestID int64           `json:"id"`
	Params    json.RawMessage `json:"params"`
}

func (r *baseRequest) Validate() error {
	const supportVersion = "2.0"
	if r.Version != supportVersion {
		return fmt.Errorf("unsuported version %s", r.Version)
	}
	if r.Method == "" {
		return fmt.Errorf("empty method passed")
	}
	if len(r.Params) == 0 {
		return fmt.Errorf("empty params")
	}
	return nil
}
