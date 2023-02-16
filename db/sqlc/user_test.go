package db

import (
	"context"
	"database/sql"
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

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomString(10)
	arg := UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	}
	newUser, err := testQuires.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, newUser.Username, oldUser.Username)
	require.Equal(t, newUser.HashedPassword, oldUser.HashedPassword)
	require.NotEqual(t, newUser.FullName, oldUser.FullName)
	require.Equal(t, newFullName, newUser.FullName)
	require.Equal(t, newUser.Email, oldUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	arg := UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	}
	newUser, err := testQuires.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, newUser.Username, oldUser.Username)
	require.Equal(t, newUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.NotEqual(t, newUser.Email, oldUser.Email)
	require.Equal(t, newEmail, newUser.Email)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashedPassword(newPassword)
	require.NoError(t, err)
	arg := UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	newUser, err := testQuires.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, newUser.Username, oldUser.Username)
	require.NotEqual(t, newUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, newUser.Email, oldUser.Email)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)
	require.WithinDuration(t, newUser.PasswordChangedAt, time.Now(), time.Second)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomString(10)
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashedPassword(newPassword)
	require.NoError(t, err)
	arg := UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	newUser, err := testQuires.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, newUser.Username, oldUser.Username)
	require.NotEqual(t, newUser.HashedPassword, oldUser.HashedPassword)
	require.NotEqual(t, newUser.FullName, oldUser.FullName)
	require.NotEqual(t, newUser.Email, oldUser.Email)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)
	require.Equal(t, newFullName, newUser.FullName)
	require.Equal(t, newEmail, newUser.Email)
	require.WithinDuration(t, newUser.PasswordChangedAt, time.Now(), time.Second)
}
