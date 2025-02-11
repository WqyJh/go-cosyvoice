package cosyvoice

const (
	wsUrl = "wss://dashscope.aliyuncs.com/api-ws/v1/inference"
)

type ClientConfig struct {
	ApiKey string `json:"api_key"`
	WsUrl  string `json:"ws_url"`
}

func DefaultConfig(apiKey string) ClientConfig {
	return ClientConfig{
		ApiKey: apiKey,
		WsUrl:  wsUrl,
	}
}
