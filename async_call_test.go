package cosyvoice_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/WqyJh/go-cosyvoice"
	"github.com/stretchr/testify/require"
)

func TestAsyncCall(t *testing.T) {
	key := os.Getenv("COSYVOICE_API_KEY")
	if key == "" {
		t.Skip("COSYVOICE_API_KEY is not set")
	}
	ctx := context.Background()
	client := cosyvoice.NewClient(key)
	asyncSynthesizer, err := client.AsyncSynthesizer(ctx)
	require.NoError(t, err)

	outputCh, err := asyncSynthesizer.RunTask(ctx)
	require.NoError(t, err)
	err = asyncSynthesizer.SendText(ctx, "你好，世界！")
	require.NoError(t, err)
	err = asyncSynthesizer.SendText(ctx, "hello, world!")
	require.NoError(t, err)
	err = asyncSynthesizer.FinishTask(ctx)
	require.NoError(t, err)

	require.Greater(t, len(outputCh), 0)
	for result := range outputCh {
		require.NoError(t, result.Err)
		require.Greater(t, len(result.Data), 0)
	}

	err = asyncSynthesizer.Close()
	require.NoError(t, err)
}

func TestAsyncPing(t *testing.T) {
	key := os.Getenv("COSYVOICE_API_KEY")
	if key == "" {
		t.Skip("COSYVOICE_API_KEY is not set")
	}
	ctx := context.Background()
	client := cosyvoice.NewClient(key)

	asyncSynthesizer, err := client.AsyncSynthesizer(ctx, cosyvoice.WithPingInterval(30))
	require.NoError(t, err)

	time.Sleep(61 * time.Second)

	outputCh, err := asyncSynthesizer.RunTask(ctx)
	require.NoError(t, err)
	err = asyncSynthesizer.SendText(ctx, "你好，世界！")
	require.NoError(t, err)

	err = asyncSynthesizer.SendText(ctx, "hello, world!")
	require.NoError(t, err)

	err = asyncSynthesizer.FinishTask(ctx)
	require.NoError(t, err)

	require.Greater(t, len(outputCh), 0)
	for result := range outputCh {
		require.NoError(t, result.Err)
		require.Greater(t, len(result.Data), 0)
	}

	err = asyncSynthesizer.Close()
	require.NoError(t, err)
}
