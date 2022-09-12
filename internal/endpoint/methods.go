package endpoint

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tmvrus/json-api/internal/entities"
)

func (e *Endpoint) getBalance(ctx context.Context, id int64, params json.RawMessage) (json.RawMessage, *apiError) {
	var req entities.BalanceRequest
	if apiErr := unmarshalAndValidate(params, &req); apiErr != nil {
		return nil, apiErr
	}

	res, err := e.balance.GetBalance(ctx, req)
	if err != nil {
		return nil, &apiError{
			Code:    codeInternalError,
			Message: "Internal storage error",
			Data:    err.Error(),
		}
	}

	return marshal(&res)
}

func (e *Endpoint) rollbackTransaction(ctx context.Context, id int64, params json.RawMessage) (json.RawMessage, *apiError) {
	var req entities.RollbackTransactionRequest
	if apiErr := unmarshalAndValidate(params, &req); apiErr != nil {
		return nil, apiErr
	}

	err := e.balance.RollbackTransaction(ctx, req.TxRef)
	if err != nil {
		return nil, &apiError{
			Code:    codeInternalError,
			Message: "Internal storage error",
			Data:    err.Error(),
		}
	}

	return nil, nil
}

func (e *Endpoint) withdrawAndDeposit(ctx context.Context, id int64, params json.RawMessage) (json.RawMessage, *apiError) {
	var req entities.WithdrawAndDepositRequest
	if apiErr := unmarshalAndValidate(params, &req); apiErr != nil {
		return nil, apiErr
	}

	res, err := e.balance.WithdrawAndDeposit(ctx, req)
	if err != nil {
		apiErr := &apiError{
			Code:    codeInternalError,
			Message: "Internal storage error",
			Data:    err.Error(),
		}
		if errors.Is(err, entities.ErrBalanceNotEnough) {
			apiErr.Code = codeBalanceNotEnough
			apiErr.Message = "Balance not enough"
			return nil, apiErr
		}
		if errors.Is(err, entities.ErrTxWasRolledBack) {
			apiErr.Code = codeTxWasRolledBack
			apiErr.Message = "transaction was rolled back"
			return nil, apiErr
		}

		return nil, apiErr
	}

	return marshal(&res)
}

func unmarshalAndValidate(data json.RawMessage, req interface{}) *apiError {
	type validator interface {
		Validate() error
	}

	if err := json.Unmarshal(data, req); err != nil {
		return &apiError{
			Code:    codeInvalidJSON,
			Message: "Invalid json",
			Data:    err.Error(),
		}
	}

	v, ok := req.(validator)
	if !ok {
		return &apiError{
			Code:    codeInternalError,
			Message: "Internal storage error",
			Data:    "Server cant validate request",
		}
	}

	if err := v.Validate(); err != nil {
		return &apiError{
			Code:    codeInvalidMethodParams,
			Message: "Invalid params json",
			Data:    err.Error(),
		}
	}

	return nil
}

func marshal(res interface{}) (json.RawMessage, *apiError) {
	data, err := json.Marshal(&res)
	if err != nil {
		return nil, &apiError{
			Code:    codeInternalError,
			Message: "Internal marshal error",
			Data:    err.Error(),
		}
	}
	return data, nil
}
