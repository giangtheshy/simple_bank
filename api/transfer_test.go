package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/giangtheshy/simple_bank/db/mock"
	db "github.com/giangtheshy/simple_bank/db/sqlc"
	"github.com/giangtheshy/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransferAPI(t *testing.T) {
	user1,_ := createRandomUser(t)
	user2,_ := createRandomUser(t)
	account1 := randomAccount(user1.Username)
	account2 := randomAccount(user2.Username)

	account1.Currency = util.USD
	account2.Currency =util.USD

	transferReq := db.TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        100,
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil).AnyTimes()
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil).AnyTimes()
	store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(transferReq)).Times(1).AnyTimes()

	server := NewTestServer(t,store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("%s/transfer",ApiPrefix)
	data, err := json.Marshal(gin.H{
		"from_account_id": account1.ID,
		"to_account_id":   account2.ID,
		"amount":          100,
		"currency":        util.USD,
	})
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
}
