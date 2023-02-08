package db

import (
	"context"
	"testing"
	"time"

	"github.com/backendmaster/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t, err)

	args := CreateUserParams{
		Username:       util.RandomOwnerName(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwnerName(),
		Email:          util.RandomEmail(),
	}

	user, err := testQuires.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	getUser, err := testQuires.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, getUser)

	require.Equal(t, getUser.Username, user.Username)
	require.Equal(t, getUser.HashedPassword, user.HashedPassword)
	require.Equal(t, getUser.FullName, user.FullName)
	require.Equal(t, getUser.Email, user.Email)

	require.WithinDuration(t, getUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
	require.WithinDuration(t, getUser.CreatedAt, user.CreatedAt, time.Second)
}
