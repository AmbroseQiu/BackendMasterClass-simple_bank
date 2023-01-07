package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/backendmaster/simple_bank/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	accessTokenDuration time.Duration,
) {
	accessToken, err := tokenMaker.CreateToken(username, accessTokenDuration)
	require.NoError(t, err)
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, accessToken)
	request.Header.Set(authorizationType, authorizationHeader)
}

func TestMiddleWare(t *testing.T) {

	testCase := []struct {
		name          string
		addAuth       func(request *http.Request, tokenMaker token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			addAuth: func(request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		// {
		// 	name: "no authorization",
		// 	addAuth: func(request *http.Request, tokenMaker token.Maker) {
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "unsupported authorizatio type",
		// 	addAuth: func(request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, "unsupported", "user", time.Minute)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "invalid authorizatio type",
		// 	addAuth: func(request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, "", "user", time.Minute)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "expired token",
		// 	addAuth: func(request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", -time.Minute)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 	},
		// },
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			routePath := "/auth"
			server.router.GET(routePath, authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, routePath, nil)
			require.NoError(t, err)

			tc.addAuth(request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
