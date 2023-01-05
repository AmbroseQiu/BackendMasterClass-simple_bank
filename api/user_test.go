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

	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"fullname": user.FullName,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
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
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				rsp := createUserResponse{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchedUser(t, recorder.Body, rsp)
			},
		},
		{
			name: "Invalid Email",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"fullname": user.FullName,
				"email":    "abcemail.com",
			},
			buildStubs: func(store *mockdb.MockStore) {
				//build stub
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Password",
			body: gin.H{
				"username": user.Username,
				"password": "12345",
				"fullname": user.FullName,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				//build stub
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Username",
			body: gin.H{
				"username": "abdc#",
				"password": password,
				"fullname": user.FullName,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				//build stub
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			//buildstubs
			tc.buildStubs(store)
			//start new server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			//marshal body to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			//check response
			tc.checkResponse(recorder)
		})
	}
}

func randomUser() (user db.User, password string) {
	return db.User{
		Username: util.RandomOwnerName(),
		FullName: util.RandomOwnerName(),
		Email:    util.RandomEmail(),
	}, util.RandomString(6)
}

func requiredBodyMatchedUser(t *testing.T, body *bytes.Buffer, rsp createUserResponse) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotResponse createUserResponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)
	require.Equal(t, rsp, gotResponse)
}
