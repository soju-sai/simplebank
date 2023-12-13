package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(7)
	hp1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hp1)

	err = CheckHashedPassword(password, hp1)
	require.NoError(t, err)

	diffPassword := RandomString(7)
	err = CheckHashedPassword(diffPassword, hp1)
	require.Error(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hp2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hp2)
	require.NotEqual(t, hp1, hp2)
}
