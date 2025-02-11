package cosyvoice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	mockApiKey := "test"
	client := NewClient(mockApiKey)
	require.Equal(t, client.config.ApiKey, mockApiKey)

	config := DefaultConfig(mockApiKey)
	client = NewClientWithConfig(config)
	require.Equal(t, mockApiKey, client.config.ApiKey)
}
