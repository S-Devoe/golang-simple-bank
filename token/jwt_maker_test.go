package token

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	secret := util.GenerateRandomString(32)
	maker, err := NewJWTMaker(secret)
	require.NoError(t, err)

	username := util.GenerateRandomUserName()
	email := util.GenerateRandomEmail()

	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	secret := util.GenerateRandomString(32)
	maker, err := NewJWTMaker(secret)
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.GenerateRandomUserName(), util.GenerateRandomEmail(), -time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidToken(t *testing.T) {
	maker, err := NewJWTMaker(util.GenerateRandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken("invalid-token")
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestInvalidAlgorithm(t *testing.T) {
	maker, err := NewJWTMaker(util.GenerateRandomString(32))
	require.NoError(t, err)

	// Create a valid token
	token, payload, err := maker.CreateToken(util.GenerateRandomUserName(), util.GenerateRandomEmail(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	// Tamper with the token to use an invalid algorithm
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3)

	// Replace the header with an invalid algorithm
	invalidHeader := base64.StdEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	parts[0] = invalidHeader
	invalidToken := strings.Join(parts, ".")

	// Verify the tampered token
	payload, err = maker.VerifyToken(invalidToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
