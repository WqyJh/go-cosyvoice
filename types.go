package cosyvoice

import (
	"encoding/json"
	"fmt"
)

type Header struct {
	Action     Action                 `json:"action"`
	TaskID     string                 `json:"task_id"`
	Streaming  string                 `json:"streaming"`
	Event      string                 `json:"event"`
	Attributes map[string]interface{} `json:"attributes"`
	Error      Error                  `json:"error,omitempty"`
}

type Error struct {
	Code    string `json:"error_code,omitempty"`
	Message string `json:"error_message,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

type Format string

const (
	PCM Format = "pcm"
	WAV Format = "wav"
	MP3 Format = "mp3"
)

type SampleRate int

const (
	SampleRate8000  SampleRate = 8000
	SampleRate16000 SampleRate = 16000
	SampleRate22050 SampleRate = 22050
	SampleRate24000 SampleRate = 24000
	SampleRate44100 SampleRate = 44100
	SampleRate48000 SampleRate = 48000
)

const (
	TextTypePlainText string = "PlainText"
)

type Model string

const (
	ModelCosyvoiceV1 Model = "cosyvoice-v1"
)

type Payload struct {
	TaskGroup  string     `json:"task_group"`
	Task       string     `json:"task"`
	Function   string     `json:"function"`
	Model      Model      `json:"model"`
	Parameters Params     `json:"parameters"`
	Resources  []Resource `json:"resources"`
	Input      Input      `json:"input"`
}

type Params struct {
	TextType   string     `json:"text_type"`        // PlainText
	Voice      Voice      `json:"voice"`            // longwan, longcheng, longhua, etc.
	Format     Format     `json:"format"`           // pcm, wav, mp3
	SampleRate SampleRate `json:"sample_rate"`      // 8000, 16000, 22050, 24000, 44100, 48000
	Volume     int        `json:"volume,omitempty"` // range: 0~100, default: 50
	Rate       int        `json:"rate,omitempty"`   // range: 0.5~2, default: 1.0
	Pitch      int        `json:"pitch,omitempty"`  // range: 0.5~2, default: 1.0
}

type Resource struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
}

type Input struct {
	Text string `json:"text"`
}

type Event struct {
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
}

type SynthesizerConfig struct {
	Model      Model      `json:"model"`            // cosyvoice-v1
	Voice      Voice      `json:"voice"`            // longwan, longcheng, longhua, etc.
	Format     Format     `json:"format"`           // pcm, wav, mp3
	SampleRate SampleRate `json:"sample_rate"`      // 8000, 16000, 22050, 24000, 44100, 48000
	Volume     int        `json:"volume,omitempty"` // range: 0~100, default: 50
	Rate       int        `json:"rate,omitempty"`   // range: 0.5~2, default: 1.0
	Pitch      int        `json:"pitch,omitempty"`  // range: 0.5~2, default: 1.0
}

func DefaultSynthesizerConfig() SynthesizerConfig {
	return SynthesizerConfig{
		Model:      "cosyvoice-v1",
		Voice:      Longxiaochun,
		Format:     "mp3",
		SampleRate: 16000,
	}
}

type Voice string

const (
	Longwan       Voice = "longwan"
	Longcheng     Voice = "longcheng"
	Longhua       Voice = "longhua"
	Longxiaochun  Voice = "longxiaochun"
	Longxiaoxia   Voice = "longxiaoxia"
	Longxiaocheng Voice = "longxiaocheng"
	Longxiaobai   Voice = "longxiaobai"
	Longlaotie    Voice = "longlaotie"
	Longshu       Voice = "longshu"
	Longjing      Voice = "longjing"
	Longmiao      Voice = "longmiao"
	Longyue       Voice = "longyue"
	Longyuan      Voice = "longyuan"
	Longfei       Voice = "longfei"
	Longjielidou  Voice = "longjielidou"
	Longshuo      Voice = "longshuo"
	Longtong      Voice = "longtong"
	Longxiang     Voice = "longxiang"
	Loongstella   Voice = "loongstella"
	Loongbella    Voice = "loongbella"
)

type Action string

const (
	ActionRunTask      Action = "run-task"
	ActionContinueTask Action = "continue-task"
	ActionFinishTask   Action = "finish-task"
)

const (
	StreamingDuplex string = "duplex"
)

const (
	TaskGroupAudio string = "audio"
)

const (
	TaskTTS string = "tts"
)

const (
	FunctionSpeechSynthesizer string = "SpeechSynthesizer"
)

func generateRunTaskCmd(taskID string, voiceConfig SynthesizerConfig) (string, error) {
	cmd := Event{
		Header: Header{
			Action:    ActionRunTask,
			TaskID:    taskID,
			Streaming: StreamingDuplex,
		},
		Payload: Payload{
			TaskGroup: TaskGroupAudio,
			Task:      TaskTTS,
			Function:  FunctionSpeechSynthesizer,
			Model:     voiceConfig.Model,
			Parameters: Params{
				TextType:   TextTypePlainText,
				Voice:      voiceConfig.Voice,
				Format:     voiceConfig.Format,
				SampleRate: voiceConfig.SampleRate,
				Volume:     voiceConfig.Volume,
				Rate:       voiceConfig.Rate,
				Pitch:      voiceConfig.Pitch,
			},
			Input: Input{},
		},
	}
	data, err := json.Marshal(cmd)
	return string(data), err
}

func generateContinueTaskCmd(taskID string, text string) (string, error) {
	cmd := Event{
		Header: Header{
			Action:    ActionContinueTask,
			TaskID:    taskID,
			Streaming: StreamingDuplex,
		},
		Payload: Payload{
			Input: Input{
				Text: text,
			},
		},
	}
	data, err := json.Marshal(cmd)
	return string(data), err
}

func generateFinishTaskCmd(taskID string) (string, error) {
	cmd := Event{
		Header: Header{
			Action:    ActionFinishTask,
			TaskID:    taskID,
			Streaming: StreamingDuplex,
		},
		Payload: Payload{
			Input: Input{},
		},
	}
	data, err := json.Marshal(cmd)
	return string(data), err
}
