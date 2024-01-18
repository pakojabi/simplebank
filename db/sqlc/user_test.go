package db

import (
	"context"
	"testing"
	"time"

	"github.com/pakojabi/simplebank/util"
	"github.com/stretchr/testify/require"
)


func TestCreateUser(t *testing.T) {
	defer cleanup()

	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	defer cleanup()

	user1 := createRandomUser(t)
	user2, err  := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second*10)
}

func createTestUser(t *testing.T, username string, hashed_password string, full_name string, email string) User {
	arg := CreateUserParams{
		Username: username,
		HashedPassword: hashed_password,
		FullName: full_name,
		Email: email,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.False(t, user.CreatedAt.IsZero())
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	return createTestUser(
		t, 
		util.RandomOwner(), 
		hashedPassword, 
		util.RandomString(6) + " " + util.RandomString(6),
		util.RandomString(4) + "@" + util.RandomString(5) + "." + util.RandomString(3),
	)
}
