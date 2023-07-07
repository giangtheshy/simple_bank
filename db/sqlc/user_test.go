package db

import (
	"context"
	"testing"
	"time"

	"github.com/giangtheshy/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user := createRandomUser(t)
	require.NotEmpty(t, user)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashPassword, user2.HashPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.ChangedPasswordAt, user2.ChangedPasswordAt, time.Second)
}

func createRandomUser(t *testing.T) User {
	hashPassword,err:=util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:   util.RandomString(6),
		HashPassword: hashPassword,
		FullName:    util.RandomOwner(),
		Email:      util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		require.NoError(t, err)
	}
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashPassword, user.HashPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	return user
}
