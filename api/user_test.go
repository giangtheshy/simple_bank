package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/giangtheshy/simple_bank/db/mock"
	db "github.com/giangtheshy/simple_bank/db/sqlc"
	"github.com/giangtheshy/simple_bank/util"
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

	err := util.CheckPasswordHash(e.password, arg.HashPassword)
	if err != nil {
		return false
	}

	e.arg.HashPassword = arg.HashPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}
func TestCreateUserAPI(t *testing.T) {
	user,password := createRandomUser(t)
	arg := db.CreateUserParams{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		HashPassword: user.HashPassword,
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg,password)).Times(1).Return(user, nil).AnyTimes()

	server := NewTestServer(t,store)
	recorder := httptest.NewRecorder()

	url  := fmt.Sprintf("%s/users",ApiPrefix)
	data,err:=json.Marshal(gin.H{
		"username" : user.Username,
		"full_name" : user.FullName,
		"email" : user.Email,
		"password" : password,
	})
	require.NoError(t, err)
	request,err:=http.NewRequest(http.MethodPost,url,bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder,request)

}

func createRandomUser(t *testing.T) (user db.User,password string){
	password = util.RandomString(6)
	hashPassword,err:=util.HashPassword(password)
	require.NoError(t, err)
	user = db.User{
		Username: util.RandomString(6),
		HashPassword: hashPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}
	return
}