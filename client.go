package cosyvoice

import (
	"context"
	"time"

	"fmt"
	"net/http"

	openairt "github.com/WqyJh/go-openai-realtime"
	gorilla "github.com/WqyJh/go-openai-realtime/contrib/ws-gorilla"
)

type Client struct {
	config ClientConfig
}

func NewClient(apiKey string) *Client {
	config := DefaultConfig(apiKey)
	return NewClientWithConfig(config)
}

func NewClientWithConfig(config ClientConfig) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) getHeaders() http.Header {
	header := make(http.Header)
	header.Add("X-DashScope-DataInspection", "enable")
	header.Add("Authorization", fmt.Sprintf("bearer %s", c.config.ApiKey))

	return header
}

type synthesizerOption struct {
	dialer            openairt.WebSocketDialer
	logger            openairt.Logger
	synthesizerConfig SynthesizerConfig
	pingInterval      int
	chanSize          int
}

type SynthesizerOption func(*synthesizerOption)

func WithDialer(dialer openairt.WebSocketDialer) SynthesizerOption {
	return func(opts *synthesizerOption) {
		opts.dialer = dialer
	}
}

func WithSynthesizerConfig(config SynthesizerConfig) SynthesizerOption {
	return func(opts *synthesizerOption) {
		opts.synthesizerConfig = config
	}
}

func WithPingInterval(interval int) SynthesizerOption {
	return func(opts *synthesizerOption) {
		opts.pingInterval = interval
	}
}

func WithChanSize(size int) SynthesizerOption {
	return func(opts *synthesizerOption) {
		opts.chanSize = size
	}
}

func WithLogger(logger openairt.Logger) SynthesizerOption {
	return func(opts *synthesizerOption) {
		opts.logger = logger
	}
}

func (c *Client) AsyncSynthesizer(ctx context.Context, opts ...SynthesizerOption) (*AsyncSynthesizer, error) {
	option := synthesizerOption{
		pingInterval:      45,
		chanSize:          32,
		dialer:            gorilla.NewWebSocketDialer(gorilla.WebSocketOptions{}),
		logger:            openairt.NopLogger{},
		synthesizerConfig: DefaultSynthesizerConfig(),
	}

	for _, opt := range opts {
		opt(&option)
	}

	header := c.getHeaders()
	url := c.config.WsUrl

	socketConn, err := option.dialer.Dial(ctx, url, header)
	if err != nil {
		return nil, err
	}

	conn := newWsConn(
		ctx,
		socketConn,
		option.logger,
		time.Duration(option.pingInterval)*time.Second,
	)

	go conn.handleHealthCheck()

	asyncSynthesizer := AsyncSynthesizer{
		conn:     conn,
		config:   option.synthesizerConfig,
		chanSize: option.chanSize,
	}

	return &asyncSynthesizer, nil
}

func (c *Client) SyncSynthesizer(ctx context.Context, opts ...SynthesizerOption) (*SyncSynthesizer, error) {
	asyncSynthesizer, err := c.AsyncSynthesizer(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &SyncSynthesizer{
		asyncSynthesizer: asyncSynthesizer,
	}, nil
}
