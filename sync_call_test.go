package cosyvoice_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/WqyJh/go-cosyvoice"
	"github.com/stretchr/testify/require"
)

func TestSyncCall(t *testing.T) {
	key := os.Getenv("COSYVOICE_API_KEY")
	if key == "" {
		t.Skip("COSYVOICE_API_KEY is not set")
	}
	ctx := context.Background()
	client := cosyvoice.NewClient(key)
	syncSynthesizer, err := client.SyncSynthesizer(ctx)
	require.NoError(t, err)

	outputCh, err := syncSynthesizer.Call(ctx, "你好，世界！")
	require.NoError(t, err)

	require.Greater(t, len(outputCh), 0)
	for result := range outputCh {
		require.NoError(t, result.Err)
		require.Greater(t, len(result.Data), 0)
	}

	outputCh, err = syncSynthesizer.Call(ctx, "hello, world!")
	require.NoError(t, err)

	require.Greater(t, len(outputCh), 0)
	for result := range outputCh {
		require.NoError(t, result.Err)
		require.Greater(t, len(result.Data), 0)
	}

	err = syncSynthesizer.Close()
	require.NoError(t, err)
}

func TestSyncPing(t *testing.T) {
	key := os.Getenv("COSYVOICE_API_KEY")
	if key == "" {
		t.Skip("COSYVOICE_API_KEY is not set")
	}
	ctx := context.Background()
	client := cosyvoice.NewClient(key)
	syncSynthesizer, err := client.SyncSynthesizer(ctx, cosyvoice.WithPingInterval(30))
	require.NoError(t, err)

	time.Sleep(61 * time.Second)

	outputCh, err := syncSynthesizer.Call(ctx, "你好，世界！")
	require.NoError(t, err)

	require.Greater(t, len(outputCh), 0)
	for result := range outputCh {
		require.NoError(t, result.Err)
		require.Greater(t, len(result.Data), 0)
	}

	time.Sleep(61 * time.Second)

	outputCh, err = syncSynthesizer.Call(ctx, "hello, world!")
	require.NoError(t, err)

	require.Greater(t, len(outputCh), 0)
	for result := range outputCh {
		require.NoError(t, result.Err)
		require.Greater(t, len(result.Data), 0)
	}

	err = syncSynthesizer.Close()
	require.NoError(t, err)
}
