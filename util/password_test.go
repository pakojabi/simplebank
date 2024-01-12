package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCorrectPassword(t *testing.T) {
	password := RandomString(6)

	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	require.NoError(t, CheckPassword(password, hashed))
}

func TestWrongPassword(t *testing.T) {
	password := RandomString(6)

	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	require.Error(t, CheckPassword("wrong", hashed))
}