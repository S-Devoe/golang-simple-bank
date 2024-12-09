package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/S-Devoe/golang-simple-bank/db/mock"
	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/S-Devoe/golang-simple-bank/util/password"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserAPI(t *testing.T) {
	user, user_password := randomUserInfo(t)
	require.NotEmpty(t, user_password)

	testCases := []struct {
		name          string
		username      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireGetUserBodyMatch(t, recorder.Body, user)
			},
		},
		{
			name:     "USER_NOT_FOUND",
			username: "randomusername",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "INTERNAL_ERROR",
			username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(user, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%s", tc.username)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// fmt.Printf("Response Body: %s\n", recorder.Body.String()) //added for debugging
			tc.checkResponse(recorder)
		})
	}
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	_, err := password.ComparePasswordAndHash(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, user_password := randomUserInfo(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  user_password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, user_password)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireCreateUserBodyMatch(t, recorder.Body, user)
			},
		},
		{
			name: "INTERNAL_ERROR",
			body: gin.H{
				"username":  user.Username,
				"password":  user_password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DUPLICATE_USERNAME",
			body: gin.H{
				"username":  user.Username,
				"password":  user_password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.
						Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "INVALID_USERNAME",
			body: gin.H{
				"username":  "in",
				"password":  user_password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "INVALID_EMAIL",
			body: gin.H{
				"username":  user.Username,
				"password":  user_password,
				"full_name": user.FullName,
				"email":     "invaliemail",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// marshal body data to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// fmt.Printf("Response Body: %s\n", recorder.Body.String()) //added for debugging
			tc.checkResponse(recorder)
		})
	}
}

func randomUserInfo(t *testing.T) (db.User, string) {
	userPassword := util.GenerateRandomString(8) // Generate random password
	hashedPassword, err := password.GeneratePasswordHash(userPassword)
	require.NoError(t, err)
	require.NotNil(t, hashedPassword)
	// Return a user with all required fields set
	return db.User{
		Username:          util.GenerateRandomString(8),
		FullName:          util.GenerateRandomName(),
		Email:             util.GenerateRandomEmail(),
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
	}, userPassword
}

func requireCreateUserBodyMatch(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response struct {
		Status int          `json:"status"`
		Data   userResponse `json:"data"`
		Error  *string      `json:"error"`
	}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, response.Status)
	require.Nil(t, response.Error)

	require.Equal(t, user.Username, response.Data.Username)
	require.Equal(t, user.FullName, response.Data.FullName)
	require.Equal(t, user.Email, response.Data.Email)

	// Validate timestamps within an acceptable range
	require.WithinDuration(t, user.CreatedAt, response.Data.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, response.Data.PasswordChanged, time.Second)

}
func requireGetUserBodyMatch(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response struct {
		Status int          `json:"status"`
		Data   userResponse `json:"data"`
		Error  *string      `json:"error"`
	}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.Status)
	require.Nil(t, response.Error)

	require.Equal(t, user.Username, response.Data.Username)
	require.Equal(t, user.FullName, response.Data.FullName)
	require.Equal(t, user.Email, response.Data.Email)

	// Validate timestamps within an acceptable range
	require.WithinDuration(t, user.CreatedAt, response.Data.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, response.Data.PasswordChanged, time.Second)

}
