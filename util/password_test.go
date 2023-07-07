package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashPassword, err := HashPassword(password)

	require.NoError(t, err)

	err = CheckPasswordHash(password, hashPassword)
	require.NoError(t, err)
}
