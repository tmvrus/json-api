package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tmvrus/json-api/internal/storage"
)

type apiMethod func(context.Context, int64, json.RawMessage) (json.RawMessage, *apiError)

type Endpoint struct {
	balance storage.BalanceStorage
	methods map[string]apiMethod
}

func NewEndpoint(s storage.BalanceStorage) *Endpoint {
	e := &Endpoint{
		balance: s,
	}
	e.methods = map[string]apiMethod{
		"getBalance":          e.getBalance,
		"rollbackTransaction": e.rollbackTransaction,
		"withdrawAndDeposit":  e.withdrawAndDeposit,
	}

	return e
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}
	if c := r.Header.Get("Content-Type"); c != "application/json" {
		http.Error(w, "Expected content type application/json, got "+c, http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body read error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var req baseRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeResponse(w, newErrResponse(0, apiError{
			Code:    codeInvalidJSON,
			Message: "invalid json",
			Data:    fmt.Sprintf("got invalid json: %q", string(body)),
		}))
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(w, newErrResponse(req.RequestID, apiError{
			Code:    codeInvalidMethodParams,
			Message: "invalid request",
			Data:    fmt.Sprintf("got invalid request: %s", err.Error()),
		}))
		return
	}

	method, ok := e.methods[req.Method]
	if !ok {
		writeResponse(w, newErrResponse(req.RequestID, apiError{
			Code:    codeMethodNotFound,
			Message: "method not found",
			Data:    fmt.Sprintf("requested rpc-method %s is not found", req.Method),
		}))
		return
	}

	res, rpcError := method(r.Context(), req.RequestID, req.Params)
	if rpcError != nil {
		writeResponse(w, newErrResponse(req.RequestID, *rpcError))
		return
	}
	writeResponse(w, newResponse(req.RequestID, res))
}
