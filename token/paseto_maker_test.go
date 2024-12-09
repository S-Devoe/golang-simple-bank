package token

import (
	"testing"
	"time"

	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	secret := util.GenerateRandomString(32)
	maker, err := NewPasetoMaker(secret)
	require.NoError(t, err)

	username := util.GenerateRandomUserName()
	email := util.GenerateRandomEmail()

	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	secret := util.GenerateRandomString(32)
	maker, err := NewPasetoMaker(secret)
	require.NoError(t, err)

	token, err := maker.CreateToken(util.GenerateRandomUserName(), util.GenerateRandomEmail(), -time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
