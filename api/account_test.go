package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/S-Devoe/golang-simple-bank/db/mock"
	db "github.com/S-Devoe/golang-simple-bank/db/sqlc"
	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateAccountAPI(t *testing.T) {
	// Mock data
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	// Mock `CreateAccount`
	store.EXPECT().
		CreateAccount(gomock.Any(), gomock.Any()).
		Times(1).
		Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	// Convert the request body to JSON
	body, err := json.Marshal(map[string]interface{}{
		"owner":    account.Owner,
		"balance":  account.Balance,
		"currency": account.Currency,
	})
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	require.Nil(t, err)

	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusCreated, recorder.Code)

	requireBodyMatchCreateAccount(t, recorder.Body, account)

}

func TestListAccountsAPI(t *testing.T) {
	// Mock data
	accounts := []db.Account{
		{ID: 1, Owner: "owner1", Balance: 100.0, Currency: "USD"},
		{ID: 2, Owner: "owner2", Balance: 200.0, Currency: "EUR"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	// Mock `ListAccounts`
	store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return(accounts, nil)
	// Mock `CountAccounts`
	store.EXPECT().
		CountAccounts(gomock.Any()).
		Times(1).
		Return(int64(len(accounts)), nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/accounts", nil)

	require.Nil(t, err)
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	requireBodyMatchAccountList(t, recorder.Body, accounts)

}

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NOTFOUND",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{
			name:      "INTERNALERROR",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.
						Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
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

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// fmt.Printf("Response Body: %s\n", recorder.Body.String()), added for debugging
			tc.checkResponse(t, recorder)
		})
	}

}

func randomUser() db.User {

	return db.User{
		Username:       util.GenerateRandomString(10),
		HashedPassword: util.GenerateRandomString(10),
		FullName:       util.GenerateRandomName(),
		Email:          util.GenerateRandomEmail(),
	}
}

func randomAccount() db.Account {
	user := randomUser()
	return db.Account{
		ID:       util.GenerateRandomInt(1, 100),
		Owner:    user.Username,
		Balance:  float64(util.GenerateRandomInt(0, 100)),
		Currency: util.GenerateRandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response struct {
		Status int        `json:"status"`
		Data   db.Account `json:"data"`
		Error  *string    `json:"error"`
	}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	// Validate the wrapped fields
	require.Equal(t, http.StatusOK, response.Status)
	require.Nil(t, response.Error)
	require.Equal(t, account, response.Data)
}

func requireBodyMatchAccountList(t *testing.T, body *bytes.Buffer, expectedAccounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response struct {
		Status int `json:"status"`
		Data   struct {
			Results []db.Account `json:"results"`
			Page    int          `json:"page"`
			Limit   int          `json:"limit"`
			Total   int          `json:"total"`
		} `json:"data"`
		Error *string `json:"error"`
	}

	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	// Validate the wrapped fields
	require.Equal(t, http.StatusOK, response.Status)
	require.Nil(t, response.Error)

	// Validate pagination metadata
	require.Equal(t, 1, response.Data.Page)   //  Expected page
	require.Equal(t, 10, response.Data.Limit) //  Expected limit
	require.NotZero(t, response.Data.Total)   //  Ensure total is not zero

	// Validate the accounts list
	require.Equal(t, len(expectedAccounts), len(response.Data.Results))
	for i, account := range expectedAccounts {
		require.Equal(t, account, response.Data.Results[i])
	}
}

func requireBodyMatchCreateAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response struct {
		Status int        `json:"status"`
		Data   db.Account `json:"data"`
		Error  *string    `json:"error"`
	}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	// Validate the wrapped fields
	require.Equal(t, http.StatusCreated, response.Status)
	require.Nil(t, response.Error)
	require.Equal(t, account, response.Data)
}
