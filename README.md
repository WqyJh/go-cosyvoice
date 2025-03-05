# CosyVoice SDK for Golang

[![Go Reference](https://pkg.go.dev/badge/github.com/WqyJh/go-cosyvoice.svg)](https://pkg.go.dev/github.com/WqyJh/go-cosyvoice)
[![Go Report Card](https://goreportcard.com/badge/github.com/WqyJh/go-cosyvoice)](https://goreportcard.com/report/github.com/WqyJh/go-cosyvoice)
[![codecov](https://codecov.io/gh/WqyJh/go-cosyvoice/branch/main/graph/badge.svg?token=bCbIfHLIsW)](https://codecov.io/gh/WqyJh/go-cosyvoice)

This library provides unofficial Go clients for [cosyvoice-websocket-api](https://help.aliyun.com/zh/model-studio/developer-reference/cosyvoice-websocket-api).

## Installation

```bash
go get github.com/WqyJh/go-cosyvoice
```

Currently, go-cosyvoice requires Go version 1.19 or greater.

## Usage


### Async Synthesizer (Streaming)


```go
import (
	"context"
	"fmt"
	"os"

	"github.com/WqyJh/go-cosyvoice"
)

func main() {
	client := cosyvoice.NewClient("your-api-key")

	ctx := context.Background()

	asyncSynthesizer, err := client.AsyncSynthesizer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer asyncSynthesizer.Close()

    output, err := synthesizer.RunTask(ctx)
	if err != nil {
		log.Println("run task err", err)
		return
	}

    texts := []string{"Hello, world!", "你好，世界！"}
    for _, text := range texts {
		err := synthesizer.SendText(ctx, text)
		if err != nil {
			log.Println("send text err", err)
			return
		}
	}

	err = synthesizer.FinishTask(ctx)
	if err != nil {
		log.Println("finish task err", err)
		return
	}

    for result := range output {
        if result.Err != nil {
            log.Println("result err", result.Err)
            return
        }
        // handle audio data
        log.Println("received audio", len(result.Data))
    }
}
```

### Sync Synthesizer (Non-Streaming)


```go
import (
	"context"
	"fmt"
	"os"

	"github.com/WqyJh/go-cosyvoice"
)

func main() {
	client := cosyvoice.NewClient("your-api-key")

	ctx := context.Background()

	syncSynthesizer, err := client.SyncSynthesizer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer syncSynthesizer.Close()

    output, err := syncSynthesizer.Call(ctx, "Hello, world!")
    if err != nil {
        log.Println("call err", err)
        return
    }

    for result := range output {
        if result.Err != nil {
            log.Println("result err", result.Err)
            return
        }
        // handle audio data
        log.Println("received audio", len(result.Data))
    }

    // another call
    output, err = syncSynthesizer.Call(ctx, "你好，世界！")
    if err != nil {
        log.Println("call err", err)
        return
    }
}
```


### Use another WebSocket dialer

The default websocket library is [coder/websocket](https://github.com/coder/websocket).
You can use another websocket library by setting a custom dialer.


```go
import (
	gorilla "github.com/WqyJh/go-cosyvoice/contrib/ws-gorilla"
)


func main() {
	syncSynthesizer, err := client.SyncSynthesizer(ctx,
		cosyvoice.WithDialer(gorilla.NewWebSocketDialer(gorilla.WebSocketOptions{})),
	)
}
```

# ADDATIONS
### The conditions for WebSocket disconnection
1. If task execution fails, the server sends the task-failed event and closes the websocketConnection. 
2. If the interval between two tasks exceeds 60 seconds, the connection will be disconnected due to timeout. In this package, after the Asynthesizer/Synthesizer object is created, a goroutine is started for health detection by send pingMessage. Developers don't need to worry about this condition. 
3. In task, if the interval between two commands exceeds 23S, the server triggers a timeout error and disconnects. When you use Asyncthesizer, if your business is likely to exceed 23S, periodically call SendText() to send empty characters or punctuation marks to prevent the connection breaking. Synthesizer is not needed.