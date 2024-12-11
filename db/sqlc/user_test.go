package db

import (
	"context"
	"testing"

	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/S-Devoe/golang-simple-bank/util/password"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user_password := util.GenerateRandomString(8)
	hashedPassword, err := password.GeneratePasswordHash(user_password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.GenerateRandomString(8),
		FullName:       util.GenerateRandomName(),
		HashedPassword: hashedPassword,
		Email:          util.GenerateRandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)

}
