package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/MasonPhan2110/SimpleBank/db/mock"
	db "github.com/MasonPhan2110/SimpleBank/db/sqlc"
	"github.com/MasonPhan2110/SimpleBank/token"
	"github.com/MasonPhan2110/SimpleBank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {

	amount := int64(10)

	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	user3, _ := randomUser(t)

	account1 := randomAccount(user1.Username)
	account2 := randomAccount(user2.Username)
	account3 := randomAccount(user3.Username)

	account1.Currency = util.VND
	account2.Currency = util.VND
	account3.Currency = util.USD

	testCases := []struct {
		name          string
		transfer      gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			transfer: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "VND",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				})).Times(1)
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
				"amount":          amount,
				"currency":        "VND",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.AssignableToTypeOf(db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account1.ID,
					Amount:        amount,
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

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()

			url := "/transfers"

			data, err := json.Marshal(tc.transfer)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
