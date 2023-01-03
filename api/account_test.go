package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/backendmaster/simple_bank/db/mock"
	db "github.com/backendmaster/simple_bank/db/sqlc"
	"github.com/backendmaster/simple_bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccout(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	//build stub
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	//start new server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	//check response
	require.Equal(t, http.StatusOK, recorder.Code)
	requiredBodyMatched(t, recorder.Body, account)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 100),
		Owner:    util.RandomOwnerName(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func requiredBodyMatched(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
