package password

import (
	"testing"

	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {

	password := util.GenerateRandomString(8)
	hashedPassword, err := GeneratePasswordHash(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	matchedPassword, err := ComparePasswordAndHash(password, hashedPassword)
	require.NoError(t, err)
	require.True(t, matchedPassword)

	wrongPassword := util.GenerateRandomString(9)
	hashedWrongPassword, err := GeneratePasswordHash(wrongPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedWrongPassword)

	matchedWrongPassword, err := ComparePasswordAndHash(password, hashedWrongPassword)
	require.NoError(t, err)
	require.False(t, matchedWrongPassword)

	// test to check if an error will occur if two users use the same password
	hashedPassword2, err := GeneratePasswordHash(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword, hashedPassword2)

	matchedPassword2, err := ComparePasswordAndHash(password, hashedPassword2)
	require.NoError(t, err)
	require.True(t, matchedPassword2)

}
