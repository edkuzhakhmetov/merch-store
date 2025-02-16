package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateJWT_Success(t *testing.T) {
	token, err := generateJWT("user")
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestValidateJWT_Success(t *testing.T) {
	username := "user"
	token, err := generateJWT(username)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	actualUsername, err := validateJWT(token)
	require.NoError(t, err)
	require.NotEmpty(t, actualUsername)
	require.Equal(t, username, actualUsername)

}

func TestValidateJWT_WithInvalidToken(t *testing.T) {
	token := "user"
	_, err := validateJWT(token)
	require.Error(t, err)
}
