package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/backendmaster/simple_bank/db/mock"
	db "github.com/backendmaster/simple_bank/db/sqlc"
	"github.com/backendmaster/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransferAPI(t *testing.T) {
	amount := int64(10)
	sameCurrencyAccount1 := randomAccount()
	sameCurrencyAccount2 := randomAccount()
	differentCurrencyAccount := randomAccount()
	sameCurrencyAccount1.Currency = util.USD
	sameCurrencyAccount2.Currency = util.USD
	differentCurrencyAccount.Currency = util.EUR

	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"from_account_id": sameCurrencyAccount1.ID,
				"to_account_id":   sameCurrencyAccount2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.TransferTxParams{
					FromAccountID: sameCurrencyAccount1.ID,
					ToAccountID:   sameCurrencyAccount2.ID,
					Amount:        amount,
				}
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(sameCurrencyAccount1.ID)).
					Times(1).
					Return(sameCurrencyAccount1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(sameCurrencyAccount2.ID)).
					Times(1).
					Return(sameCurrencyAccount2, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid json body",
			body: gin.H{
				"from_account_id": -1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Currency Mismatched of Account1",
			body: gin.H{
				"from_account_id": differentCurrencyAccount.ID,
				"to_account_id":   sameCurrencyAccount2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(differentCurrencyAccount.ID)).
					Times(1).
					Return(differentCurrencyAccount, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Currency Mismatched of Account2",
			body: gin.H{
				"from_account_id": sameCurrencyAccount1.ID,
				"to_account_id":   differentCurrencyAccount.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(sameCurrencyAccount1.ID)).
					Times(1).
					Return(sameCurrencyAccount1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(differentCurrencyAccount.ID)).
					Times(1).
					Return(differentCurrencyAccount, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Account1 Not Found",
			body: gin.H{
				"from_account_id": sameCurrencyAccount1.ID,
				"to_account_id":   sameCurrencyAccount2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(sameCurrencyAccount1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Transcation Failed",
			body: gin.H{
				"from_account_id": sameCurrencyAccount1.ID,
				"to_account_id":   sameCurrencyAccount2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.TransferTxParams{
					FromAccountID: sameCurrencyAccount1.ID,
					ToAccountID:   sameCurrencyAccount2.ID,
					Amount:        amount,
				}
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(sameCurrencyAccount1.ID)).
					Times(1).
					Return(sameCurrencyAccount1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(sameCurrencyAccount2.ID)).
					Times(1).
					Return(sameCurrencyAccount2, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			//build stub
			tc.buildStubs(store)
			//start new server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			//marshal body to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			//check response
			tc.checkResponse(t, recorder)
		})
	}

}
