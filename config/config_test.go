package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJWT_Success(t *testing.T) {
	err := os.Setenv("JWT_SECRET_KEY", "key")
	require.NoError(t, err)
	err = validateJWT()
	require.NoError(t, err)
}

func TestJWT_ErrorWithEmptyKey(t *testing.T) {
	err := os.Setenv("JWT_SECRET_KEY", "")
	require.NoError(t, err)
	err = validateJWT()
	require.EqualError(t, err, "JWT secret (JWT_SECRET_KEY) is not set")
}

func TestValidateDatabase_Success(t *testing.T) {
	err := os.Setenv("DB_HOST", "host")
	require.NoError(t, err)

	err = os.Setenv("DB_PORT", "1111")
	require.NoError(t, err)

	err = os.Setenv("DB_USER", "user")
	require.NoError(t, err)

	err = validateDatabase()
	require.NoError(t, err)
}

func TestValidateDatabase_ErrorWithEmptyHost(t *testing.T) {
	err := os.Setenv("DB_HOST", "")
	require.NoError(t, err)

	err = os.Setenv("DB_PORT", "1111")
	require.NoError(t, err)

	err = os.Setenv("DB_USER", "user")
	require.NoError(t, err)

	err = os.Setenv("DB_PASS", "pas")
	require.NoError(t, err)

	err = os.Setenv("DB_NAME", "name")
	require.NoError(t, err)

	err = os.Setenv("DB_SSLMODE", "disable")
	require.NoError(t, err)

	err = validateDatabase()
	require.EqualError(t, err, "database host address (DB_HOST) is not set")
}

func TestValidateDatabase_ErrorWithEmptyPort(t *testing.T) {
	err := os.Setenv("DB_HOST", "host")
	require.NoError(t, err)

	err = os.Setenv("DB_PORT", "")
	require.NoError(t, err)

	err = os.Setenv("DB_USER", "user")
	require.NoError(t, err)

	err = os.Setenv("DB_PASS", "pas")
	require.NoError(t, err)

	err = os.Setenv("DB_NAME", "name")
	require.NoError(t, err)

	err = os.Setenv("DB_SSLMODE", "disable")
	require.NoError(t, err)

	err = validateDatabase()
	require.EqualError(t, err, "database port (DB_PORT) is not set")
}

func TestValidateDatabase_ErrorWithEmptyUser(t *testing.T) {
	err := os.Setenv("DB_HOST", "host")
	require.NoError(t, err)

	err = os.Setenv("DB_PORT", "1111")
	require.NoError(t, err)

	err = os.Setenv("DB_USER", "")
	require.NoError(t, err)

	err = os.Setenv("DB_PASS", "pas")
	require.NoError(t, err)

	err = os.Setenv("DB_NAME", "name")
	require.NoError(t, err)

	err = os.Setenv("DB_SSLMODE", "disable")
	require.NoError(t, err)

	err = validateDatabase()
	require.EqualError(t, err, "database username (DB_USER) is not set")
}
