package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/MasonPhan2110/SimpleBank/db/mock"
	db "github.com/MasonPhan2110/SimpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	account1 := db.Account{
		ID:       1,
		Owner:    "minh",
		Balance:  1000,
		Currency: "VND",
	}
	account2 := db.Account{
		ID:       2,
		Owner:    "minh",
		Balance:  0,
		Currency: "VND",
	}
	testCases := []struct {
		name          string
		transfer      gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			transfer: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          10,
				"currency":        "VND",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.AssignableToTypeOf(db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        10,
				})).Times(1).Return(db.TransferTxResult{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid ID",
			transfer: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account1.ID,
				"amount":          10,
				"currency":        "VND",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.AssignableToTypeOf(db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account1.ID,
					Amount:        10,
				})).Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				fmt.Println(recorder.Body)
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

			// start test server and send request

			server := NewServer(store)

			recorder := httptest.NewRecorder()

			url := "/transfers"

			data, err := json.Marshal(tc.transfer)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
