package main

import (
	"context"
	"log"
	"os"

	openairt "github.com/WqyJh/go-openai-realtime"
	pcm "github.com/WqyJh/go-openai-realtime/examples/voice/pcm"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"

	"github.com/WqyJh/go-cosyvoice"
)

const (
	sampleRate = cosyvoice.SampleRate44100
)

func main() {
	key := os.Getenv("COSYVOICE_API_KEY")
	if key == "" {
		log.Println("COSYVOICE_API_KEY is not set")
		return
	}

	client := cosyvoice.NewClient(key)

	ctx := context.Background()

	// // Init speaker
	speaker.Init(beep.SampleRate(sampleRate), int(sampleRate)/10)
	defer speaker.Close()
	streamers := make(chan *pcm.PCMStream, 100)
	beep.NewBuffer(beep.Format{
		SampleRate:  beep.SampleRate(sampleRate),
		NumChannels: 1,
		Precision:   2,
	})
	playDone := make(chan bool)
	go func() {
		defer close(playDone)
		done := make(chan bool)
		speaker.Play(beep.Iterate(func() beep.Streamer {
			stream, ok := <-streamers
			if !ok {
				close(done)
				return nil
			}
			return stream
		}))
		<-done
	}()

	synthsizerConfig := cosyvoice.SynthesizerConfig{
		Model:      cosyvoice.ModelCosyvoiceV1,
		Voice:      cosyvoice.Longshu,
		Format:     cosyvoice.FormatPCM,
		SampleRate: sampleRate,
	}
	synthesizer, err := client.SyncSynthesizer(ctx,
		cosyvoice.WithSynthesizerConfig(synthsizerConfig),
		cosyvoice.WithLogger(openairt.StdLogger{}),
	)
	if err != nil {
		log.Println("websocket conn err", err)
		return
	}

	text := "你好，今天的天气怎么样?"
	OnceCall(synthesizer, ctx, text, streamers)

	texts1 := "挺不错的。"
	OnceCall(synthesizer, ctx, texts1, streamers)

	close(streamers)
	// wait audio playback completed
	<-playDone

	err = synthesizer.Close()
	if err != nil {
		log.Println(err)
	}
}

func OnceCall(synthesizer *cosyvoice.SyncSynthesizer, ctx context.Context, text string, streamers chan *pcm.PCMStream) {
	resultCh, err := synthesizer.Call(ctx, text)
	if err != nil {
		log.Println("call err", err)
		return
	}
	for result := range resultCh {
		if result.Err != nil {
			log.Println("result err", result.Err)
			return
		}
		log.Println("received audio", len(result.Data))
		if len(result.Data) > 0 {
			streamer := pcm.NewPCMStream(result.Data, beep.Format{SampleRate: beep.SampleRate(sampleRate), NumChannels: 1, Precision: 2})
			streamers <- streamer
		}
	}
}
