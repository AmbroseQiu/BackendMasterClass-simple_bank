package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/backendmaster/simple_bank/db/mock"
	db "github.com/backendmaster/simple_bank/db/sqlc"
	"github.com/backendmaster/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)

	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("match params %v and password %v", e.arg, e.password)
}

func EqCreateUserParamsMatcher(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	arg := db.CreateUserParams{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}
	//build stub
	store.EXPECT().
		CreateUser(gomock.Any(), EqCreateUserParamsMatcher(arg, password)).
		Times(1).
		Return(user, nil)
	//start new server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	body := gin.H{
		"username": user.Username,
		"password": password,
		"fullname": user.FullName,
		"email":    user.Email,
	}
	//marshal body to json
	data, err := json.Marshal(body)
	require.NoError(t, err)
	url := "/users"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
	//check response
	// tc.checkResponse(t, recorder)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func randomUser() (user db.User, password string) {
	return db.User{
		Username: util.RandomOwnerName(),
		FullName: util.RandomOwnerName(),
		Email:    util.RandomEmail(),
	}, util.RandomString(6)
}

func requiredBodyMatchedUser(t *testing.T, body *bytes.Buffer, rsp createUserRequest) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotResponse createUserResponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)
	require.Equal(t, rsp, gotResponse)
}
