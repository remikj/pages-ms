package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewServerFromEnv_shouldInitializeServerFromEnv_whenEnvsValid(t *testing.T) {
	t.Setenv("SERVICE_PORT", "8081")
	config, err := ConfigurationFromEnv()

	assert.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, 8081, config.Port)
}

func TestNewServerFromEnv_shouldInitializeDefaultServerFromEnv_whenEnvsNotSet(t *testing.T) {
	require.NoError(t, os.Unsetenv("SERVICE_PORT"))

	config, err := ConfigurationFromEnv()

	assert.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, 8080, config.Port)
}

func TestNewServerFromEnv_shouldReturnErr_whenEnvsInvalid(t *testing.T) {
	t.Setenv("SERVICE_PORT", "asd")
	server, err := ConfigurationFromEnv()

	assert.Nil(t, server)
	require.Error(t, err)
}
