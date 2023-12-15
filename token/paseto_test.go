package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	issuer := username
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	myClaims, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, myClaims)
	require.NotZero(t, myClaims.ID)
	require.Equal(t, issuer, myClaims.Issuer)
	require.Contains(t, "local", myClaims.Subject)

	require.WithinDuration(t, issuedAt, myClaims.IssuedAt.Time, time.Second)
	require.WithinDuration(t, issuedAt, myClaims.NotBefore.Time, time.Second)
	require.WithinDuration(t, expiredAt, myClaims.ExpiresAt.Time, time.Second)
}

func TestPasetoTokenExpired(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, claims)
}
