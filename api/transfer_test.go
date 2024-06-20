package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mockdb "github.com/kuthumipepple/bank-backend/db/mock"
	db "github.com/kuthumipepple/bank-backend/db/sqlc"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTransferAPI(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()

	account2.Currency = account1.Currency
	body := gin.H{
		"from_account_id": account1.ID,
		"to_account_id":   account2.ID,
		"amount":          10,
		"currency":        account1.Currency,
	}

	arg := db.TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}
	result := db.TransferTxResult{
		Transfer: db.Transfer{
			Amount: arg.Amount,
		},
	}

	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
		Times(1).
		Return(account1, nil)
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
		Times(1).
		Return(account2, nil)
	store.EXPECT().
		TransferTx(gomock.Any(), gomock.Eq(arg)).
		Times(1).
		Return(result, nil)

	url := "/transfers"
	data, err := json.Marshal(body)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	server := NewServer(store)

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchTransfer(t, recorder.Body, result)
}

func requireBodyMatchTransfer(t *testing.T, body *bytes.Buffer, result db.TransferTxResult) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTransfer db.TransferTxResult
	err = json.Unmarshal(data, &gotTransfer)
	require.NoError(t, err)

	require.Equal(t, result, gotTransfer)
}
