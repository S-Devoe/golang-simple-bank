package db

import (
	"context"
	"testing"

	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		Username:       util.GenerateRandomString(8),
		FullName:       util.GenerateRandomName(),
		HashedPassword: util.GenerateRandomString(8),
		Email:          util.GenerateRandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)

}
