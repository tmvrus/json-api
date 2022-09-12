package endpoint_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	api "github.com/tmvrus/json-api/internal/endpoint"
	"github.com/tmvrus/json-api/internal/endpoint/mock"
	"github.com/tmvrus/json-api/internal/entities"
)

//go:generate mockgen -destination=./mock/storage.go -package=mock -source=../storage/storage.go

func TestEndpoint_OK(t *testing.T) {
	t.Parallel()

	_ = require.Error
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	storage := mock.NewMockBalanceStorage(ctrl)
	endpoint := api.NewEndpoint(storage)

	tests := []struct {
		name      string
		in        []byte
		setupMock func(s *mock.MockBalanceStorage)
		expected  string
	}{
		{
			name: "deposit",
			in:   []byte(`{"jsonrpc":"2.0","method":"withdrawAndDeposit","id":42,"params":{"callerId":15,"playerName":"Player15","currency":"USD","transactionRef":"tx1","deposit":100}}`),
			setupMock: func(s *mock.MockBalanceStorage) {
				s.
					EXPECT().
					WithdrawAndDeposit(
						gomock.Any(),
						entities.WithdrawAndDepositRequest{
							CommonRequest: entities.CommonRequest{
								CallerID:   15,
								PlayerName: "Player15",
							},
							Currency: "USD",
							TxRef:    "tx1",
							Withdraw: 0,
							Deposit:  100,
						}).
					Return(
						entities.WithdrawAndDepositResponse{
							NewBalance: 100,
							TxID:       "tx1",
						},
						nil)
			},
			expected: `{"id":42, "jsonrpc":"2.0", "result":{"newBalance":100, "transactionId":"tx1"}}`,
		},
		{
			name: "balance",
			in:   []byte(`{"jsonrpc":"2.0","method":"getBalance","id":43,"params":{"callerId":15,"playerName":"Player15","currency":"USD"}}`),
			setupMock: func(s *mock.MockBalanceStorage) {
				s.
					EXPECT().
					GetBalance(gomock.Any(), entities.BalanceRequest{
						CommonRequest: entities.CommonRequest{
							CallerID:   15,
							PlayerName: "Player15",
						},
						Currency: "USD",
					}).Return(entities.BalanceResponse{Balance: 100}, nil)
			},
			expected: `{"id":43, "jsonrpc":"2.0", "result":{"balance":100}}`,
		},
		{
			name: "rollback",
			in:   []byte(`{"jsonrpc":"2.0","method":"rollbackTransaction","id":44,"params":{"callerId":15,"playerName":"Player15","transactionRef":"tx1"}}`),
			setupMock: func(s *mock.MockBalanceStorage) {
				s.
					EXPECT().
					RollbackTransaction(gomock.Any(), "tx1").
					Return(nil)
			},
			expected: `{"id":44, "jsonrpc":"2.0", "result":null}`,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			test.setupMock(storage)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(test.in))
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			endpoint.ServeHTTP(w, req)

			response, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)
			require.JSONEqf(t, test.expected, string(response), string(response))
		})
	}

}
