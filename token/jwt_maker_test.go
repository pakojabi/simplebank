package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pakojabi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
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

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
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

func TestInvalidJWTTokenAlgNone(t *testing.T) {

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, NewJWTPayloadClaims(payload))
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err = maker.Verify(token)
	require.Error(t, err)
	require.Equal(t, ErrInvalidToken, err)
	require.Nil(t, payload)
}
