package token

import (
	"testing"
	"time"

	"github.com/pakojabi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.Make(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.Verify(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker.Make(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	_, err = maker.Verify(token)
	require.Error(t, err)
	require.Equal(t, ErrExpiredToken, err)
}
