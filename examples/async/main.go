package main

import (
	"context"
	"log"
	"os"

	"github.com/WqyJh/go-cosyvoice"
	openairt "github.com/WqyJh/go-openai-realtime"
	pcm "github.com/WqyJh/go-openai-realtime/examples/voice/pcm"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
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

	// Init speaker
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
	synthesizer, err := client.AsyncSynthesizer(ctx,
		cosyvoice.WithSynthesizerConfig(synthsizerConfig),
		cosyvoice.WithLogger(openairt.StdLogger{}),
	)
	if err != nil {
		log.Println("websocket conn err", err)
		return
	}
	defer synthesizer.Close()

	texts := []string{"你好，今天的天气怎么样?", "武汉今天多云转晴。"}
	OnceCall(synthesizer, ctx, texts, streamers)

	texts1 := []string{"出去散散步，", "挺不错的。"}
	OnceCall(synthesizer, ctx, texts1, streamers)

	close(streamers)
	// wait audio playback completed
	<-playDone
}

func OnceCall(synthesizer *cosyvoice.AsyncSynthesizer, ctx context.Context, texts []string, streamers chan *pcm.PCMStream) {

	wait := make(chan struct{})

	output, err := synthesizer.RunTask(ctx)
	if err != nil {
		log.Println("run task err", err)
		return
	}

	go func() {
		defer close(wait)

		for result := range output {
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

	}()

	for _, text := range texts {
		log.Println("send synthesizer text", text)
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

	<-wait
}
