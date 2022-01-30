package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createdUser := createRandomUser(t)
	queriedUser, err := testQueries.GetUser(context.Background(), createdUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, queriedUser)
	require.Equal(t, createdUser.CreatedAt, queriedUser.CreatedAt)
	require.Equal(t, createdUser.Email, queriedUser.Email)
	require.Equal(t, createdUser.FullName, queriedUser.FullName)
	require.Equal(t, createdUser.HashedPassword, queriedUser.HashedPassword)
	require.Equal(t, createdUser.PasswordChangedAt, queriedUser.PasswordChangedAt)
	require.Equal(t, createdUser.Username, queriedUser.Username)
}

func createRandomUser(t *testing.T) User {
	hashedPass, err := HashPassword(randomString(8))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       randomString(8),
		HashedPassword: hashedPass,
		FullName:       randomString(6),
		Email:          randomString(6) + "@" + randomString(4) + randomString(3),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Username, user.Username)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}
