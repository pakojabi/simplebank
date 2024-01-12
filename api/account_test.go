package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
	mockdb "github.com/pakojabi/simplebank/db/mock"
	db "github.com/pakojabi/simplebank/db/sqlc"
	"github.com/pakojabi/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct{
		name string
		accountID int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, errors.New("some error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
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
		
			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()
		
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
		
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}


func TestCreateAccount(t *testing.T) {
	testCases := []struct{
		name string
		owner string
		currency string
		buildStubs func(store *mockdb.MockStore, expectedOwner, expectedCurrency string)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, expectedOwner, expectedCurrency string)
	}{
		{
			name: "OK",
			owner: "owner",
			currency: "EUR",
			buildStubs: func(store *mockdb.MockStore, expectedOwner, expectedCurrency string) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
						Owner: expectedOwner,
						Currency: expectedCurrency,
						Balance: 0,
					})).
					Times(1).
					Return(db.Account{ID: util.RandomInt(1,1000), Owner: expectedOwner, Currency: expectedCurrency, Balance: 0}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedOwner, expectedCurrency string){
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccountProps(t, recorder.Body, expectedOwner, expectedCurrency, 0)
			},
		},
		{
			name: "BadRequest",
			owner: "owner",
			currency: "YEN",
			buildStubs: func(store *mockdb.MockStore, expectedOwner, expectedCurrency string) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedOwner, expectedCurrency string){
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			owner: "owner",
			currency: "EUR",
			buildStubs: func(store *mockdb.MockStore, expectedOwner, expectedCurrency string) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedOwner, expectedCurrency string){
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unexisting owner",
			owner: "owner",
			currency: "EUR",
			buildStubs: func(store *mockdb.MockStore, expectedOwner, expectedCurrency string) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, &pq.Error{Code: pq.ErrorCode("23503")})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedOwner, expectedCurrency string){
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Duplicate account",
			owner: "owner",
			currency: "EUR",
			buildStubs: func(store *mockdb.MockStore, expectedOwner, expectedCurrency string) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, &pq.Error{Code: pq.ErrorCode("23505")})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedOwner, expectedCurrency string){
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store, tc.owner, tc.currency)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, getReaderFor(t, createAccountRequest{
				Owner: tc.owner,
				Currency: tc.currency,
			}))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder, tc.owner, tc.currency)
		})
	}
}

func TestListAccounts(t *testing.T) {
	
	var dummy_accounts []db.Account 
	for i := 0; i < 10; i++ {
		dummy_accounts = append(dummy_accounts, db.Account{
			ID: util.RandomInt(1, 1000),
			Balance: util.RandomMoney(),
			Owner: fmt.Sprintf("owner_%d", i),
			Currency: util.RandomCurrency(),
		})
	}

	testCases := []struct{
		name string
		pageId, pageSize int
		buildStubs func(store *mockdb.MockStore, expectedOffset, expectedLimit int64)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, expectedOffset, expectedLimit int64)
	}{
		{
			name: "OK",
			pageId: 1,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore, expectedOffset, expectedLimit int64) {
				store.EXPECT().
					ListAccounts(gomock.Any(), db.ListAccountsParams{
						Limit: expectedLimit,
						Offset: expectedOffset,
					}).
					Times(1).
					Return(dummy_accounts[expectedOffset:expectedLimit], nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedOffset, expectedLimit int64) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, dummy_accounts[expectedOffset:expectedLimit])
			},

		},
		{
			name: "BadRequest",
			pageId: 0,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore, expectedOffset, expectedLimit int64) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedOffset, expectedLimit int64) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},

		},
		{
			name: "InternalServerError",
			pageId: 1,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore, expectedOffset, expectedLimit int64) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedOffset, expectedLimit int64) {
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
			limit := int64(tc.pageSize)
			offset := (int64(tc.pageId - 1) * int64(tc.pageSize))
			tc.buildStubs(store, offset, limit)
			
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", tc.pageId, tc.pageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder, offset, limit)
		})
	}
}

func TestUpdateAccount(t *testing.T)  {
	account := randomAccount()

	testCases := []struct{
		name string
		accountId int64
		newBalance int64
		buildStubs func(store *mockdb.MockStore, newBalance int64)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, newBalance int64)

	}{
		{
			name: "OK",
			accountId: account.ID,
			newBalance: 100,
			buildStubs: func(store *mockdb.MockStore, newBalance int64) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), db.UpdateAccountParams{
						ID: account.ID,
						Balance: newBalance,
					}).
					Times(1).
					Return(db.Account{
						ID: account.ID,
						Owner: account.Owner,
						Balance: newBalance,
						Currency: account.Currency,
						CreatedAt: account.CreatedAt,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, newBalance int64) {
				require.Equal(t, http.StatusOK, recorder.Code)
				expectedAccount := account
				expectedAccount.Balance = newBalance
				requireBodyMatchAccount(t, recorder.Body, expectedAccount)
			},
		},
		{
			name: "BadRequest on url",
			accountId: 0,
			newBalance: 100,
			buildStubs: func(store *mockdb.MockStore, newBalance int64) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, newBalance int64) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest on body",
			accountId: account.ID,
			newBalance: 0,
			buildStubs: func(store *mockdb.MockStore, newBalance int64) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, newBalance int64) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			accountId: account.ID,
			newBalance: 100,
			buildStubs: func(store *mockdb.MockStore, newBalance int64) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, newBalance int64) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Not Found",
			accountId: account.ID,
			newBalance: 100,
			buildStubs: func(store *mockdb.MockStore, newBalance int64) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, newBalance int64) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store, tc.newBalance)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountId)
			request, err := http.NewRequest(http.MethodPut, url, getReaderFor(t, updateAccountBody{
				NewBalance: tc.newBalance,
			}))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder, tc.newBalance)
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct{
		name string
		accountID int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), account.ID).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			accountID: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), account.ID).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			server.router.ServeHTTP(recorder, request)
			
			require.NoError(t, err)
			tc.checkResponse(t, recorder)

		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID: util.RandomInt(1, 1000),
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccountProps(t *testing.T, body *bytes.Buffer, expectedOwner, expectedCurrency string, expectedBalance int64){
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, gotAccount.Balance)
	require.Equal(t, expectedOwner, gotAccount.Owner)
	require.Equal(t, expectedBalance, gotAccount.Balance)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, expectedAccounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, len(expectedAccounts), len(gotAccounts))
	for i := range gotAccounts {
		require.Equal(t, expectedAccounts[i], gotAccounts[i])
	}
}
